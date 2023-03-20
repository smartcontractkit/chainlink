import json
import os
import random
import string
from pathlib import Path
from jinja2 import Environment, FileSystemLoader
from solcx import compile_standard
from solc_ast import process_ast

def parse_events(abi_json, event_map, struct_map, file_name):
    parsed_abi = json.loads(abi_json)

    for item in parsed_abi:
        if item['type'] == 'event':
            event_key = f"{file_name}_{item['name']}"
            event_map[event_key] = item
        elif item['type'] == 'tuple':
            struct_name = item['name']
            struct_map[struct_name] = item

    return event_map, struct_map

def generate_mock_contract(events, structs):
    env = Environment(loader=FileSystemLoader('.'))
    template = env.get_template('mock_contract_template.j2')

    return template.render(events=events.values(), structs=structs.values())

def get_abi_files(root):
    return [str(p) for p in Path(root).rglob('*.abi')]

def main():
    root = "/path/to/your/abis"  # Change this to the path containing your ABI files
    abi_files = get_abi_files(root)

    event_map = {}
    struct_map = {}

    for abi_file in abi_files:
        with open(abi_file, 'r') as f:
            abi_json = f.read()

        file_name = Path(abi_file).stem
        event_map, struct_map = parse_events(abi_json, event_map, struct_map, file_name)

    # Generate the mock contract
    mock_contract = generate_mock_contract(event_map, struct_map)

    # Save the mock contract to a file
    with open("/path/to/output/EventsMock.sol", "w") as f:
        f.write(mock_contract)

    print("Generated EventsMock.sol mock contract!")

if __name__ == "__main__":
    main()
