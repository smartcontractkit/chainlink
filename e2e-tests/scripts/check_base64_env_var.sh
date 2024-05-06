#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
WHITE='\033[0;37m'
NC='\033[0m' # No Color

echo
if [ ! -z "$BASE64_CONFIG_OVERRIDE" ]; then
    echo "${GREEN}BASE64_CONFIG_OVERRIDE is set, it will be used.${NC}"
    echo "Decoded content:"
    echo
    echo "$BASE64_CONFIG_OVERRIDE" | base64 --decode
    echo
    echo
    echo "${GREEN}Press RETURN to confirm and continue...${NC} ${RED}or CTRL+C to exit${NC}"
    read -r -n 1 key
else
    echo "${YELLOW}BASE64_CONFIG_OVERRIDE is not set, checking for overrides.toml file...${NC}"
    if [ -f "testconfig/overrides.toml" ]; then
        echo "${GREEN}Found testconfig/overrides.toml file. Here's its content:${NC}"
        echo
        cat "testconfig/overrides.toml"
        echo
        echo
        echo "${GREEN}Press RETURN to base64-encode it and run the test...${NC} ${RED}or CTRL+C to exit${NC}"
        read -r -n 1 key
        if [[ $key = "" ]]; then
            export BASE64_CONFIG_OVERRIDE=$(cat testconfig/overrides.toml | base64)
        fi
    else
        echo "${RED}testconfig/overrides.toml file does not exist. Please create it or set BASE64_CONFIG_OVERRIDE manually.${NC}"
        echo
    fi
fi