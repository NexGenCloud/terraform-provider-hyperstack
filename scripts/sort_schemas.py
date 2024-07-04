#!/usr/bin/env python3.11

"""
Script Name: sort_schemas.py
Description:
    Normalizes API json to have standard representation

Usage:
    Script expects single input parameter with API
    specification JSON.

    To update API specification in place:
    $ python sort_schemas.py "api.json"
"""

import argparse
import json
import canonicaljson

def fix_api_spec(api_file: str) -> None:
  """
  Updates API file in place.

  Args:
      api_file: Path to schema file
  """
  with open(api_file, 'r') as file:
    data = json.load(file)

  data = canonicaljson.encode_canonical_json(data)
  data = json.loads(data.decode('utf-8'))

  with open(api_file, 'w') as file:
    json.dump(data, file, indent=4)


if __name__ == "__main__":
  parser = argparse.ArgumentParser(
    description='Updates JSON to a normalized state',
  )
  parser.add_argument(
    'api_file',
    type=str,
    help='Path to the API spec file',
  )
  args = parser.parse_args()

  fix_api_spec(args.api_file)
