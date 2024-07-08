#!/bin/bash

# Check if folder argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <folder>"
  exit 1
fi

# Get the folder name from the command argument
folder=$1

# Create the uml/$folder directory
uml_folder="uml/$folder"
mkdir -p "$uml_folder"

# Walk through all subfolders of ./$folder
find "$folder" -type f -name '*.sol' | while read -r file; do
   # Get the relative path of the file
  relative_path="${file#$folder/}"
 
  # Create the corresponding directory structure in uml/$folder
  uml_subfolder="$uml_folder/$(dirname "$relative_path")"
  mkdir -p "$uml_subfolder"
  
  # Flatten the Solidity file
  flattened_file="flattened_$(basename "$file")"
  forge flatten "$file" -o "$flattened_file"
  
  # Generate UML diagram
  sol2uml "$flattened_file" -o "$uml_subfolder"
  
  # Clean up the flattened file
  rm "$flattened_file"
done

echo "UML diagrams generated in $uml_folder"

