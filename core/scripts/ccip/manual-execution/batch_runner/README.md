# batchrun.py

This script will run the manual-execution script repeatedly using a CSV file to
supply message IDs and txn hashes.

## Usage

Build main.go and copy 'manual-execution' into the same directory as this script.

Export messages to a file named 'msgs.csv' in the same directory as this script.
CSV file will ignore the first line (header) and each line should have two values separated by a comma.

The first value is the message ID and the second value is the CCIP send transaction hash.

Example csv file:
```
"message_id","transaction_hash"
0x1e221l7db3f193d19353d42e1bcece771e7edar57149a7f30afd31r8aa783e9a,0x03c88qfd30ar54f36a353262a67362838r52029a47h2a31141a0430bda937ba2
0x1e221l7db3f193d19353d42e1bcece771e7edar57149a7f30afd31r8aa783e9a,0x03c88qfd30ar54f36a353262a67362838r52029a47h2a31141a0430bda937ba2
0x1e221l7db3f193d19353d42e1bcece771e7edar57149a7f30afd31r8aa783e9a,0x03c88qfd30ar54f36a353262a67362838r52029a47h2a31141a0430bda937ba2
0x1e221l7db3f193d19353d42e1bcece771e7edar57149a7f30afd31r8aa783e9a,0x03c88qfd30ar54f36a353262a67362838r52029a47h2a31141a0430bda937ba2
```
