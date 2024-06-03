#!/usr/bin/env python3

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


def process_emptry_attr_types(json_data):
  paths = json_data.get("paths", {})
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


def process_ref_strings(json_data):
  """ Recursively process $ref strings in the JSON data. """
  if isinstance(json_data, dict):
    for key, value in list(json_data.items()):
      if key == "$ref" and isinstance(value, str):
        # Replace spaces after 'schemas/' in $ref strings
        json_data[key] = re.sub(r'\s+', '', value)
      else:
        process_ref_strings(value)
  elif isinstance(json_data, list):
    for item in json_data:
      process_ref_strings(item)


def process_components_schemas(json_data):
  """ Process the components.schemas part of the JSON to replace keys with spaces. """
  paths = json_data.get("paths", {})
  components = json_data.get("components", {})
  schemas = components.get("schemas", {})

  for key in list(schemas.keys()):
    new_key = re.sub(r'\s+', '', key)
    if new_key != key:
      schemas[new_key] = schemas.pop(key)

    props = schemas[new_key]["properties"]
    if "status" in props and "message" in props:
      print("Fixing %s" % new_key)
      del props["status"]
      del props["message"]

      # TODO: recheck
      if new_key == "Instances" and "instance_count" in props:
        del props["instance_count"]

      if len(props.keys()) > 1:
        print("Warning: check this key")

      props = {} if len(props.keys()) == 0 else list(props.values())[0]
      if "type" not in props:
        props["type"] = "object"

      schemas[new_key] = props
      # print(props)
      # print(schemas[new_key]["properties"])

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

  paths["/core/virtual-machines/{virtual_machine_id}/sg-rules"] = paths["/core/virtual-machines/{id}/sg-rules"]
  del paths["/core/virtual-machines/{id}/sg-rules"]
  paths["/core/virtual-machines/{virtual_machine_id}/sg-rules"]["post"]["parameters"][0]["name"] = "virtual_machine_id"

  paths["/core/virtual-machines/{virtual_machine_id}/sg-rules/{id}"] = paths[
    "/core/virtual-machines/{virtual_machine_id}/sg-rules/{sg_rule_id}"]
  del paths["/core/virtual-machines/{virtual_machine_id}/sg-rules/{sg_rule_id}"]
  paths["/core/virtual-machines/{virtual_machine_id}/sg-rules/{id}"]["delete"]["parameters"][0][
    "name"] = "virtual_machine_id"
  paths["/core/virtual-machines/{virtual_machine_id}/sg-rules/{id}"]["delete"]["parameters"][1]["name"] = "id"


def main(file_path):
  with open(file_path, 'r') as file:
    data = json.load(file)

  # Process $ref strings
  process_ref_strings(data)

  # Process components.schemas to replace keys with spaces
  process_components_schemas(data)

  process_emptry_attr_types(data)

  # Write the modified data back to the file
  with open(file_path, 'w') as file:
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

  main(args.spec_file)
