#!/usr/bin/env python3
import argparse
import json


# https://github.com/hashicorp/terraform-plugin-codegen-framework/issues/115
def process_components_schemas(json_data):
  """ Process the components.schemas part of the JSON to replace keys with spaces. """
  resources = json_data.get("resources", [])

  for resource in resources:
    if resource["name"] == "core_virtual_machine":
      for attr in resource["schema"]["attributes"]:
        taint_all(attr)
      resource["schema"]["attributes"] = [x for x in resource["schema"]["attributes"] if x["name"] != "count"]

def taint_all(attr):
  modifier_str = [
    {
      "custom": {
        "imports": [
          {
            "path": "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
          }
        ],
        "schema_definition": "stringplanmodifier.RequiresReplace()"
      }
    }
  ]
  modifier_bool = [
    {
      "custom": {
        "imports": [
          {
            "path": "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
          }
        ],
        "schema_definition": "boolplanmodifier.RequiresReplace()"
      }
    }
  ]
  modifier_int64 = [
    {
      "custom": {
        "imports": [
          {
            "path": "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
          }
        ],
        "schema_definition": "int64planmodifier.RequiresReplace()"
      }
    }
  ]
  modifier_object = [
    {
      "custom": {
        "imports": [
          {
            "path": "github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
          }
        ],
        "schema_definition": "objectplanmodifier.RequiresReplace()"
      }
    }
  ]

  if "string" in attr:
    attr["string"]["plan_modifiers"] = modifier_str
  if "bool" in attr:
    attr["bool"]["plan_modifiers"] = modifier_bool
  if "int64" in attr:
    attr["int64"]["plan_modifiers"] = modifier_int64
  if "single_nested" in attr:
    attr["single_nested"]["plan_modifiers"] = modifier_object
    for nested in attr["single_nested"]["attributes"]:
      taint_all(nested)

def main(file_path):
  with open(file_path, 'r') as file:
    data = json.load(file)

  # Process components.schemas to replace keys with spaces
  process_components_schemas(data)

  # Write the modified data back to the file
  with open(file_path, 'w') as file:
    json.dump(data, file, indent=4)


if __name__ == "__main__":
  parser = argparse.ArgumentParser(description='Remove spaces from schema names in a JSON file.')
  parser.add_argument('file_path', type=str, help='Path to the JSON file')
  args = parser.parse_args()

  main(args.file_path)
