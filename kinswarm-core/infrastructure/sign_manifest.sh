#!/bin/bash
# Anchoring the 10M Batch to The Keeper's Identity

ROOT="6a99f4d0755e0ce9dba8afb3f3bde5c0c23a364ad47e886ebbaeca8ba75914b2"
FILE="manifest.txt"

echo "Anchoring Root: $ROOT"
echo "$ROOT" > $FILE

# Detached armor signature using your Jon S. (Pray4Love1) key
gpg --detach-sign --armor --yes $FILE

echo "--- Sovereign Manifest Created ---"
ls -l manifest.txt*
gpg --verify manifest.txt.asc manifest.txt
