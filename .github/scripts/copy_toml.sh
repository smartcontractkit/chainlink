#!/bin/bash          
          # This script searches for TOML files matching given regex patterns and copies them to a target directory.

          echo "patterns: $1"
          echo "search_folder: $2"
          echo "target_folder: $3"

          if [ "$#" -ne 3 ]; then
              echo "Usage: $0 <pattern-array> <search-folder> <target-folder>"
              exit 1
          fi

          declare -a patterns=($1)   # Array of regex filename patterns
          search_folder=$2           # Folder to search in
          target_folder=$3           # Target folder for copying files

          if [ ! -d "$search_folder" ]; then
              echo "Search folder does not exist: $search_folder"
              exit 1
          fiwhy 

          mkdir -p "$target_folder"

          copy_files() {
              local pattern=$1
              # Find .toml files, extract filenames, and match pattern
              find "$search_folder" -type f -name "*.toml" | while read -r file; do
                  filename=$(basename "$file")
                  if [[ $filename =~ $pattern ]]; then
                      cp "$file" "$target_folder"
                  fi
              done
          }

          for pattern in "${patterns[@]}"; do
              copy_files "$pattern"
          done
