#!/bin/bash

go run main.go \
 deploy-contracts \
 --ocrfile=ocr_config.json \
 --chainid=11155111 \
 --onlysetconfig=false \
 --ethurl=#<ETH_URL_HERE> \
 --accountkey=#<ACCOUNT_KEY_HERE> 
