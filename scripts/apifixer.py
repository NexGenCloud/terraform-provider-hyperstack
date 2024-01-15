#!/usr/bin/env python3
import argparse
import json
import re


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
  components = json_data.get("components", {})
  schemas = components.get("schemas", {})

  for key in list(schemas.keys()):
    new_key = re.sub(r'\s+', '', key)
    if new_key != key:
      schemas[new_key] = schemas.pop(key)


def main(file_path):
  with open(file_path, 'r') as file:
    data = json.load(file)

  # Process $ref strings
  process_ref_strings(data)

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
