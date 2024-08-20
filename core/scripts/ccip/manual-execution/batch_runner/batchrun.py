#!/usr/bin/env python3

import subprocess
import os
import tempfile
import json

binary_name = "./manual-execution"
input_msgs = "msgs.csv"

# global config
src_rpc = "<rpc for source chain here>"
dest_rpc = "<rpc for destination chain here>"
dest_owner_key = "<private key for destination chain here>"
commit_store = "<commit store address here>"
off_ramp = "<off ramp address here>"
on_ramp = "<on ramp address here>"
dest_start_block = 11063581 # set to block where the messages commit report was written.
gas_limit_override = 2000000 # you can change this or leave it as is.

lines = open(input_msgs, "r").read().split("\n")[1:]

for i, pairs in enumerate(lines):
    # per msg config
    parts = pairs.split(",")
    if len(parts) != 2:
        if pairs != "":
            print("skipping CSV line with unexpected format: %s" % pairs)
        continue

    msg_id = parts[0]
    ccip_send_tx = parts[1]

    print("[%d/%d] >>> %s %s" % (i, len(lines), ccip_send_tx, msg_id))

    config = {
        "source_chain_tx": ccip_send_tx,
        "ccip_msg_id": msg_id,
        "src_rpc": src_rpc,
        "dest_rpc": dest_rpc,
        "dest_owner_key": dest_owner_key,
        "commit_store": commit_store,
        "off_ramp": off_ramp,
        "dest_start_block": dest_start_block,
        "gas_limit_override": gas_limit_override
    }
    json_config = json.dumps(config)

    with open("config.json", 'w') as f:
        f.write(json_config)

    try:
        subprocess.run([binary_name])
    except subprocess.CalledProcessError as e:
        print("called process error: ", e)
