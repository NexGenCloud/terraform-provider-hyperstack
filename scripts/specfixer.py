#!/usr/bin/env python3
import argparse
import json


# https://github.com/hashicorp/terraform-plugin-codegen-framework/issues/115
def process_components_schemas(json_data):
  """ Process the components.schemas part of the JSON to replace keys with spaces. """
  resources = json_data.get("resources", [])

  for resource in resources:
    if resource["name"] == "core_virtual_machine_sg_rule":
      for attr in resource["schema"]["attributes"]:
        taint_all(attr)

    if resource["name"] == "core_virtual_machine":
      resource["schema"]["blocks"] = []

      for attr in resource["schema"]["attributes"]:
        match attr["name"]:
          case "profile":
            del attr["single_nested"]["computed_optional_required"]
            resource["schema"]["blocks"].append(attr)
            for row in attr["single_nested"]["attributes"]:
              if row["name"] == "name":
                # TODO: fix this
                # https://github.com/hashicorp/terraform-plugin-framework/issues/740
                # https://discuss.hashicorp.com/t/optional-block-with-required-attribute-in-framework/54371
                row["string"]["computed_optional_required"] = "optional"
            # https://github.com/hashicorp/terraform-plugin-framework/issues/603
            attr["list_nested"] = attr["single_nested"]
            del attr["single_nested"]
            attr["list_nested"]["nested_object"] = {
              "attributes": attr["list_nested"]["attributes"]
            }
            del attr["list_nested"]["attributes"]
          case "assign_floating_ip":
            attr["bool"]["default"] = {
              "static": False,
            }
          case "create_bootable_volume":
            attr["bool"]["default"] = {
              "static": False,
            }

        if attr["name"] not in ["flavor", "environment", "image", "keypair", "created_at", "floating_ip_status", "id", "status"]:
          taint_all(attr)

      resource["schema"]["attributes"] = [x for x in resource["schema"]["attributes"] if x["name"] not in ["count", "profile"]]


def taint_all(attr, toskip=[]):
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
