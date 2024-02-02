#!/bin/bash

# Function to convert a comma-separated list into a TOML array format.
# Usage: convert_to_toml_array "elem1,elem2,elem3"
# Effect: "a,b,c" -> ["a","b","c"]
function convert_to_toml_array() {
    local IFS=','
    local input_array=($1)
    local toml_array_format="["

    for element in "${input_array[@]}"; do
        toml_array_format+="\"$element\","
    done

    toml_array_format="${toml_array_format%,}]"
    echo "$toml_array_format"
}