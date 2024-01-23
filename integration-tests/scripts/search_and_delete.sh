#!/bin/bash

if [ $# -eq 0 ]; then
    read -p "Enter the comma-separated list of filenames to search for: " filenames
else
    filenames=$@
fi

IFS=',' read -ra filenames_arr <<< "$filenames"

# Start search from the current working directory
current_dir=$(pwd)

echo "Searching for files in $current_dir..."
found_files=()
for filename in "${filenames_arr[@]}"; do
    while IFS= read -r file; do
        found_files+=("$file")
    done < <(find "$current_dir" -type f -name "$filename" 2>/dev/null)
done

if [[ ${#found_files[@]} -eq 0 ]]; then
    echo "No files found."
    exit 0
fi

echo "Found files:"
for file in "${found_files[@]}"; do
    echo "$file"
done

read -p "Do you want to remove all these files? (y/n): " confirm

if [[ $confirm == "yes" ]] || [[ $confirm == "y" ]]; then
    for file in "${found_files[@]}"; do
        rm "$file"
    done
    echo "Files removed."
else
    echo "Files not removed."
fi