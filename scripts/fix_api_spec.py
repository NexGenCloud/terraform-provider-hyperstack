#!/usr/bin/env python3.11

"""
Script Name: fix_api_spec.py
Description:
    Modifies OpenAPI Hyperstack specification to fix various issues with
    fields due to GET/POST/PUT merging and API endpoint inconsistencies.

    Specification:
    https://swagger.io/specification/

Usage:
    Script expects single input parameter with API
    specification JSON.

    To update specification in place:
    $ python fix_api_spec.py "api.json"

Notes:
    Ideally this file should be continuously updated to reflect latest
    API state, and reducing logic here is crucial
"""

import argparse
import json
import re
import copy
from typing import Dict, Any

AttrType = Dict[str, Any]


def attr_remove_ref_spaces(data: AttrType) -> None:
  """
  Iteratively goes through the document removing any spaces from
  $ref references in JSON schema.

  Args:
      data: Data chunk to process.
  """
  if isinstance(data, dict):
    for key, value in list(data.items()):
      if key == "$ref" and isinstance(value, str):
        # Replace spaces after 'schemas/' in $ref strings
        data[key] = re.sub(r'(\s+|%20)', '', value)
      else:
        attr_remove_ref_spaces(value)
  elif isinstance(data, list):
    for item in data:
      attr_remove_ref_spaces(item)


def attr_fix_empty_types(data: AttrType) -> None:
  """
  Replaces all empty schema types to strings.

  Args:
      data: Data chunk to process.
  """
  paths = data.get("paths", {})
  for path in paths:
    methods = paths[path]
    for method in methods:
      if "parameters" in methods[method]:
        for param in methods[method]["parameters"]:
          if "schema" in param:
            # if there is a key [type] and it is empty, set it to string
            if "type" not in param["schema"] or param["schema"]["type"] == "":
              param["schema"]["type"] = "string"
              print("Fixing empty attribute type in %s" % path)


def attr_fix_components(data: AttrType) -> None:
  """
  Goes through all schemas applying various fixes to API definitions.

  Args:
      data: Data chunk to process.
  """
  paths = data.get("paths", {})
  components = data.get("components", {})
  schemas = components.get("schemas", {})

  # Pass 1: normalize schema names (strip stray spaces / %20 in keys).
  for key in list(schemas.keys()):
    new_key = re.sub(r'(\s+|%20)', '', key)
    if new_key != key:
      schemas[new_key] = schemas.pop(key)

  # Pass 2: schema-specific fixes below. These must run before the envelope
  # unwrapping (Pass 3) because they rely on the wrapped structure still being
  # present (e.g. the RBAC role fix reads properties off the wrapper).

  # API incorrectly returns "roles"
  schemas["RbacRoleDetailResponseModelFixed"] = copy.deepcopy(schemas["RbacRoleDetailResponseModel"])
  schemas["RbacRoleDetailResponseModelFixed"]["properties"]["roles"] = schemas["RbacRoleDetailResponseModelFixed"]["properties"]["role"]
  del schemas["RbacRoleDetailResponseModelFixed"]["properties"]["role"]
  paths["/auth/roles/{id}"]["get"]["responses"]["200"]["content"]["application/json"]["schema"]["$ref"] = "#/components/schemas/RbacRoleDetailResponseModelFixed"

  schemas["Flavor_Fields"]["properties"]["ram"]["type"] = "number"

  schemas["Instance_Overview_Fields"]["properties"]["ram"]["type"] = "number"
  schemas["Container_Overview_Fields"]["properties"]["ram"]["type"] = "number"
  schemas["Instance_Flavor_Fields"]["properties"]["ram"]["type"] = "number"
  schemas["Cluster_Flavor_Fields"]["properties"]["ram"]["type"] = "number"

  # TODO: UNSYNCED see tf provider
  #schemas["ImportKeypairPayload"]["properties"]["environment"] = schemas["ImportKeypairPayload"]["properties"]["environment_name"]
  #del schemas["ImportKeypairPayload"]["properties"]["environment_name"]

  # TODO: UNSYNCED see tf provider
  schemas["Create_Security_Rule_Payload"]["required"].append("virtual_machine_id")
  # schemas["Create_Security_Rule_Payload"]["properties"]["virtual_machine_id"] = {
  #   "type": "integer",
  # }
  schemas["Create_Security_Rule_Payload"]["properties"]["port_range_min"] = {
    "type": "integer",
  }
  schemas["Create_Security_Rule_Payload"]["properties"]["port_range_max"] = {
    "type": "integer",
  }

  schemas["ClusterFields"]["properties"]["kube_config"] = {
    "type": "string"
  }

  # GPU stock configuration keys are digit-prefixed (1x, 2x, ... 10x), which
  # all sanitize to the same invalid Go/Terraform identifier ("x") and clash.
  # Prefix them with "n" so each maps to a distinct attribute (n1x, n2x, ...),
  # matching the SDK's N-prefixed fields used by the stocks data source.
  if "New_Configurations_Response" in schemas:
    config_props = schemas["New_Configurations_Response"].get("properties", {})
    for p in list(config_props.keys()):
      config_props["n%s" % p] = config_props.pop(p)

  # Pass 3: unwrap envelope response bodies.
  unwrap_envelopes(schemas)

  # Pass 4: normalize paths (parameter names + synthesized sg-rules read).
  normalize_paths(paths, schemas)


def unwrap_envelopes(schemas: AttrType) -> None:
  """
  Collapses API "envelope" response schemas down to their payload.

  Most list/detail responses wrap the real payload under a single field
  alongside envelope noise (status/message) and pagination noise
  (page/page_size/count), e.g. ``{"status", "message", "volume": {...}}``.
  Left as-is these surface as extra schema attributes that either duplicate
  identically-named query parameters (page_size) or nest the real fields one
  level too deep (breaking the dot-paths used by ``schema.ignores`` in
  generator-config.yml). We strip the noise and, when a single payload field
  remains, replace the schema with that field so downstream schemas stay flat.

  Args:
      schemas: The components/schemas map, mutated in place.
  """
  noise = ["status", "message", "page", "page_size", "count"]
  for name in list(schemas.keys()):
    schema = schemas[name]
    if not isinstance(schema, dict):
      continue
    props = schema.get("properties")
    if not props or "status" not in props or "message" not in props:
      continue

    print("Fixing %s" % name)
    for key in noise:
      props.pop(key, None)

    # Instances list carries an extra count field under a different name
    if name == "Instances" and "instance_count" in props:
      del props["instance_count"]

    if len(props) == 0:
      schemas[name] = {"type": "object"}
    elif len(props) == 1:
      payload = list(props.values())[0]
      if "type" not in payload and "$ref" not in payload:
        payload["type"] = "object"
      schemas[name] = payload
    # More than one payload field remains: keep it as an object with the
    # noise removed rather than guessing which field to unwrap.


def rename_path_param(operation: AttrType, old_name: str, new_name: str) -> None:
  """
  Renames a single path parameter within an operation in place.

  Args:
      operation: OpenAPI operation object (e.g. paths[path]["get"]).
      old_name: Current parameter name.
      new_name: Desired parameter name.
  """
  for param in operation.get("parameters", []):
    if param.get("in") == "path" and param.get("name") == old_name:
      param["name"] = new_name


def normalize_paths(paths: AttrType, schemas: AttrType) -> None:
  """
  Normalizes API paths to the stable parameter names expected by
  generator-config.yml and reconstructs the sg-rules read endpoint.

  Upstream renamed path parameters (`{id}` -> `{volume_id}`,
  `{virtual_machine_id}`/`{id}` -> `{vm_id}`/`{sg_rule_id}`) and dropped the
  sg-rules listing endpoint. generator-config.yml derives the
  `virtual_machine_id`/`id` schema attributes from these path parameter
  names, so we normalize them back here instead of churning the config.
  """
  # Volume detail: {volume_id} -> {id}. The GET response is wrapped in a
  # Volume envelope; point it at the flat Volume_Fields so the resource
  # schema exposes the volume attributes directly.
  vol_old = "/core/volumes/{volume_id}"
  vol_new = "/core/volumes/{id}"
  if vol_old in paths:
    paths[vol_new] = paths.pop(vol_old)
    for method in paths[vol_new].values():
      rename_path_param(method, "volume_id", "id")
    if "Volume_Fields" in schemas:
      paths[vol_new]["get"]["responses"]["200"]["content"]["application/json"]["schema"] = {
        "$ref": "#/components/schemas/Volume_Fields",
      }

  # SG rules collection: {vm_id} -> {virtual_machine_id}. Upstream removed the
  # GET listing, but the resource still needs a read source for its computed
  # fields, so synthesize one returning the flat Security_Group_Rule_Fields.
  sg_old = "/core/virtual-machines/{vm_id}/sg-rules"
  sg_new = "/core/virtual-machines/{virtual_machine_id}/sg-rules"
  if sg_old in paths:
    paths[sg_new] = paths.pop(sg_old)
    rename_path_param(paths[sg_new].get("post", {}), "vm_id", "virtual_machine_id")
    paths[sg_new]["get"] = {
      "parameters": [{
        "in": "path",
        "name": "virtual_machine_id",
        "required": True,
        "schema": {"type": "integer"},
      }],
      "responses": {"200": {
        "description": "Success",
        "content": {"application/json": {
          "schema": {"$ref": "#/components/schemas/Security_Group_Rule_Fields"},
        }},
      }},
    }

  # SG rules detail (DELETE): {vm_id}/{sg_rule_id} -> {virtual_machine_id}/{id}
  sg_detail_old = "/core/virtual-machines/{vm_id}/sg-rules/{sg_rule_id}"
  sg_detail_new = "/core/virtual-machines/{virtual_machine_id}/sg-rules/{id}"
  if sg_detail_old in paths:
    paths[sg_detail_new] = paths.pop(sg_detail_old)
    rename_path_param(paths[sg_detail_new].get("delete", {}), "vm_id", "virtual_machine_id")
    rename_path_param(paths[sg_detail_new].get("delete", {}), "sg_rule_id", "id")


def fix_api_spec(spec_file: str) -> None:
  """
  Updates specification file in place, applying various schema fixes.

  Args:
      spec_file: Path to schema file
  """
  with open(spec_file, 'r') as file:
    data = json.load(file)

  attr_remove_ref_spaces(data)
  attr_fix_components(data)
  attr_fix_empty_types(data)

  with open(spec_file, 'w') as file:
    json.dump(data, file, indent=4)


if __name__ == "__main__":
  parser = argparse.ArgumentParser(
    description='Fixes API specification for Nexgen Hyperstack',
  )
  parser.add_argument(
    'spec_file',
    type=str,
    help='Path to the JSON file',
  )
  args = parser.parse_args()

  fix_api_spec(args.spec_file)