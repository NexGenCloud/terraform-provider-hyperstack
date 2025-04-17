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

  for key in list(schemas.keys()):
    new_key = re.sub(r'(\s+|%20)', '', key)
    if new_key != key:
      schemas[new_key] = schemas.pop(key)

    props = schemas[new_key]["properties"]
    # TODO: how do we get it? e.g. count
    to_delete = ["status", "message", "page", "page_size", "count"]
    if "status" in props and "message" in props:
      print("Fixing %s" % new_key)

      for r in to_delete:
        if r in props:
          del props[r]


      # TODO: recheck
      if new_key == "Instances" and "instance_count" in props:
        del props["instance_count"]

      if len(props.keys()) > 1:
        pass

      subfield = True
      if len(props.keys()) == 0:
        props = {}
      elif len(props.keys()) == 1:
        props = list(props.values())[0]
      else:
        # TODO: not anymore as now we have some filters like page
        print("Warning: check this key")

      if "type" not in props and "$ref" not in props:
        props["type"] = "object"

      schemas[new_key] = props
      #print(props)
      # print(schemas[new_key]["properties"])

  schemas["InstanceOverviewFields"]["properties"]["ram"]["type"] = "number"
  schemas["ContainerOverviewFields"]["properties"]["ram"]["type"] = "number"
  schemas["InstanceFlavorFields"]["properties"]["ram"]["type"] = "number"
  schemas["FlavorFields"]["properties"]["ram"]["type"] = "number"

  schemas["ImportKeypairPayload"]["properties"]["environment"] = schemas["ImportKeypairPayload"]["properties"][
    "environment_name"]
  del schemas["ImportKeypairPayload"]["properties"]["environment_name"]

  schemas["CreateSecurityRulePayload"]["required"].append("virtual_machine_id")
  schemas["CreateSecurityRulePayload"]["properties"]["virtual_machine_id"] = {
    "type": "integer",
  }
  schemas["CreateSecurityRulePayload"]["properties"]["port_range_min"] = {
    "type": "integer",
  }
  schemas["CreateSecurityRulePayload"]["properties"]["port_range_max"] = {
    "type": "integer",
  }

  schemas["ClusterFields"]["properties"]["kube_config"] = {
    "type": "string"
  }

  # Fix digit-prefixed keys
  props = schemas["NewConfigurationsResponse"]["properties"]
  for p in list(props.keys()):
    props["N%s" % p] = props.pop(p)

  paths["/core/virtual-machines/{virtual_machine_id}/sg-rules"] = paths["/core/virtual-machines/{vm_id}/sg-rules"]
  del paths["/core/virtual-machines/{vm_id}/sg-rules"]
  paths["/core/virtual-machines/{virtual_machine_id}/sg-rules"]["post"]["parameters"][0]["name"] = "virtual_machine_id"

  paths["/core/virtual-machines/{virtual_machine_id}/sg-rules/{id}"] = paths[
    "/core/virtual-machines/{vm_id}/sg-rules/{sg_rule_id}"]
  del paths["/core/virtual-machines/{vm_id}/sg-rules/{sg_rule_id}"]
  paths["/core/virtual-machines/{virtual_machine_id}/sg-rules/{id}"]["delete"]["parameters"][0][
    "name"] = "virtual_machine_id"
  paths["/core/virtual-machines/{virtual_machine_id}/sg-rules/{id}"]["delete"]["parameters"][1]["name"] = "id"

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
