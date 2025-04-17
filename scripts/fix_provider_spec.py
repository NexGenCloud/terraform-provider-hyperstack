#!/usr/bin/env python3.11

"""
Script Name: fix_provider_spec.py
Description:
    Modifies specification generated with tfplugingen-openapi to fix
    various issues with fields due to GET/POST/PUT merging and
    API endpoint inconsistencies:
    https://github.com/hashicorp/terraform-plugin-codegen-framework/issues/115

    Specification:
    https://developer.hashicorp.com/terraform/plugin/code-generation/specification

Usage:
    Script expects single input parameter with provider
    specification JSON.

    To update specification in place:
    $ python fix_provider_spec.py "provider-spec.json"

Notes:
    Ideally this file should be continuously updated to reflect latest
    API state, and reducing logic here is crucial
"""

import argparse
import json
from typing import Dict, Callable, Any

AttrType = Dict[str, Any]


def attr_update_nested(attr: AttrType, updater: Callable[[str, AttrType, str], None]) -> None:
  """
  Iteratively goes through attribute and children (if present) applying
  updater function each time.

  Args:
      attr: Attribute to process.
      updater:
        Function to apply to each attribute. It gets attribute type,
        current attribute we are working with, and "fixed" attribute
        for objects (single_nested -> object) due to code inconsistencies
        across Terraform Plugin Framework
  """
  # https://github.com/hashicorp/terraform-plugin-framework/tree/main/resource/schema
  basic_types = [
    "string",
    "list",
    "bool",
    "int64",
    "number",
  ]
  for t in basic_types:
    if t in attr:
      updater(t, attr[t], t)
      return

  # For nested definitions we need to go through all nested attributes
  if "single_nested" in attr:
    updater("single_nested", attr["single_nested"], "object")
    
    # Make sure that in attr single_nested object has attributes array
    if attr.get("single_nested", {}).get("attributes"):
      for nested in attr["single_nested"]["attributes"]:
        attr_update_nested(nested, updater)
  if "list_nested" in attr:
    updater("list_nested", attr["list_nested"], "list")
    for nested in attr["list_nested"]["nested_object"]["attributes"]:
      attr_update_nested(nested, updater)


def attr_set_modifier(attr: AttrType, modifier: str | None) -> None:
  """
  Iteratively goes through attribute and children (if present) applying
  plan modifiers to each.

  Args:
      attr: Attribute to process.
      modifier:
        Modifier name to put in go call, e.g. RequiresReplace.
        If None, deletes any existing plan modifiers.
  """
  def mod(attr_type, attr_value, attr_subtype):
    if modifier is None:
      del attr_value["plan_modifiers"]
      return

    attr_value["plan_modifiers"] = [{
      "custom": {
        "imports": [
          {
            "path": "github.com/hashicorp/terraform-plugin-framework/resource/schema/%splanmodifier" % attr_subtype
          }
        ],
        "schema_definition": "%splanmodifier.%s()" % (attr_subtype, modifier)
      }
    }]

  attr_update_nested(attr, mod)


def attr_set_default(attr: AttrType, value: str | None) -> None:
  """
  Iteratively goes through attribute and children (if present) setting
  default values to each.

  Specification:
  https://developer.hashicorp.com/terraform/plugin/framework/resources/default

  Args:
      attr: Attribute to process.
      value:
        Default value to put in go call, should be compatible with go types.
        If None, deletes any existing plan modifiers.
  """

  def mod(attr_type, attr_value, attr_subtype):
    if value is None:
      del attr_value["default"]
      return

    attr_value["default"] = {
      "custom": {
        "imports": [
          {
            "path": "github.com/hashicorp/terraform-plugin-framework/resource/schema/%sdefault" % attr_subtype
          }
        ],
        "schema_definition": "%sdefault.StaticValue(%s)" % (attr_subtype, value)
      }
    }

  attr_update_nested(attr, mod)


def attr_update_behavior(attr: AttrType, behavior: str | None) -> None:
  """
  Iteratively goes through attribute and children (if present) setting
  schema behavior to each.

  Specification:
  https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors

  Args:
      attr: Attribute to process.
      behavior:
        Behavior to set (Optional/Required/Computed or combination).
        If None, deletes any existing plan modifiers.
  """

  def mod(attr_type, attr_value, attr_subtype):
    if behavior is None:
      del attr_value["computed_optional_required"]
      return

    attr_value["computed_optional_required"] = behavior

  attr_update_nested(attr, mod)


def fix_provider_spec(spec_file: str) -> None:
  """
  Updates specification file in place, applying various schema fixes.

  Args:
      spec_file: Path to schema file
  """
  with open(spec_file, 'r') as file:
    data = json.load(file)

  datasources = data.get("datasources", [])
  for row in datasources:
    if row["name"] == "core_keypair":
      for attr in row["schema"]["attributes"]:
        if attr["name"] == "core_keypair":
          row["schema"]["attributes"] = attr["set_nested"]["nested_object"]["attributes"]
          break
      for attr in row["schema"]["attributes"]:
        if attr["name"] == "id":
          attr["int64"]["computed_optional_required"] = "required"

  resources = data.get("resources", [])
  for row in resources:
    if row["name"] == "core_virtual_machine_sg_rule":
      for attr in row["schema"]["attributes"]:
        attr_set_modifier(attr, "RequiresReplace")

    if row["name"] == "core_cluster":
      for attr in row["schema"]["attributes"]:
        immutable_params = [
          "name",
          "environment_name",
          "kubernetes_version",
          "node_count",
          "keypair_name",
          "enable_public_ip",
        ]
        if attr["name"] in immutable_params:
          attr_set_modifier(attr, "RequiresReplace")
        computed_params = [
          "id",
          "created_at",
          "api_address",
          "kube_config",
          "status",
          "status_reason",
          "node_addresses",
        ]
        if attr["name"] in computed_params:
          attr_update_behavior(attr, "computed")

      # TODO: remove when docs are fixed
      remove_params = [
        "node_addresses",
      ]
      row["schema"]["attributes"] = [x for x in row["schema"]["attributes"] if x["name"] not in remove_params]

    if row["name"] == "core_volume":
      for attr in row["schema"]["attributes"]:
        immutable_params = [
          "name",
          "environment_name",
          "description",
          "volume_type",
          "size",
          "image_id",
          "callback_url",
        ]
        if attr["name"] in immutable_params:
          attr_set_modifier(attr, "RequiresReplace")

    if row["name"] == "core_virtual_machine":
      row["schema"]["blocks"] = []

      for attr in row["schema"]["attributes"]:
        match attr["name"]:
          case "profile":
            attr_update_behavior(attr, "computed")
            del attr["single_nested"]["computed_optional_required"]
            row["schema"]["blocks"].append(attr)
            # for attr_row in attr["single_nested"]["attributes"]:
            #   if attr_row["name"] == "name":
            #     # TODO: fix this
            #     # https://github.com/hashicorp/terraform-plugin-framework/issues/740
            #     # https://discuss.hashicorp.com/t/optional-block-with-required-attribute-in-framework/54371
            #     attr_row["string"]["computed_optional_required"] = "optional"
            # https://github.com/hashicorp/terraform-plugin-framework/issues/603
            attr["list_nested"] = attr["single_nested"]
            del attr["single_nested"]
            attr["list_nested"]["nested_object"] = {
              "attributes": attr["list_nested"]["attributes"]
            }
            del attr["list_nested"]["attributes"]
          case "flavor":
            attr_update_behavior(attr, "computed")
          case "profile":
            attr_update_behavior(attr, "computed")
          case "labels":
            # Set default value for labels to empty list
            attr_set_default(attr, "types.ListValueMust(types.StringType, []attr.Value{})")
          case "assign_floating_ip":
            attr["bool"]["default"] = {
              "static": False,
            }
          case "create_bootable_volume":
            attr["bool"]["default"] = {
              "static": False,
            }

        vm_immutable_params = [
          "assign_floating_ip",
          "count",
          "create_bootable_volume",
          "environment_name",
          "flavor_name",
          "image_name",
          "key_name",
          "name",
          "profile",
          "user_data",
          "volume_attachments",
        ]
        if attr["name"] in vm_immutable_params:
          attr_set_modifier(attr, "RequiresReplace")

        vm_mutable_params = [
          "callback_url",
        ]
        if attr["name"] in vm_mutable_params:
          attr_update_behavior(attr, "optional")

        vm_computed_params = [
          "created_at",
          "fixed_ip",
          "floating_ip",
          "floating_ip_status",
          "id",
          "locked",
          "os",
          "power_state",
          "status",
          "vm_state",
          "security_rules",
        ]
        if attr["name"] in vm_computed_params:
          attr_update_behavior(attr, "computed")
        vm_temporary_wrong_params = ["user_data"]
        if attr["name"] in vm_computed_params + vm_mutable_params + vm_temporary_wrong_params:
          attr_set_modifier(attr, "UseStateForUnknown")

      vm_remove_params = [
        "count",
        "profile",
        "contract_id",
      ]
      row["schema"]["attributes"] = [x for x in row["schema"]["attributes"] if x["name"] not in vm_remove_params]

  with open(spec_file, 'w') as file:
    json.dump(data, file, indent=4)


if __name__ == "__main__":
  parser = argparse.ArgumentParser(
    description='Fixes terraform provider specification for Nexgen Hyperstack',
  )
  parser.add_argument(
    'spec_file',
    type=str,
    help='Path to the provider spec file',
  )
  args = parser.parse_args()

  fix_provider_spec(args.spec_file)
