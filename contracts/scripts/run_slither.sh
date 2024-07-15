#!/bin/bash

# Input: Comma-separated list of files
FILE_LIST=$1

# Convert the comma-separated list into an array
IFS=',' read -r -a files <<< "$FILE_LIST"

# Loop through each file and check permissions, then run slither
for file in "${files[@]}"
do
   echo "Processing file: $file"

   # Check if file exists
   if [ ! -f "$file" ]; then
     echo "File does not exist: $file"
     continue
   fi

   # Check read permissions
   if [ ! -r "$file" ]; then
     echo "File is not readable: $file"
     continue
   fi

   # Run slither on the file
   echo "Running slither on $file"
   slither "$file" || echo "Slither failed on $file"
done
