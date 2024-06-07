#!/bin/bash

go run main.go \
 deploy-contracts \
 --ocrfile=ocr_config.json \
 --chainid=11155111 \
 --ethurl=ETH_URL \
 --accountkey=ACCOUNT_KEY \
 --onlysetconfig=false \
 --skipfunding=false \
 --dryrun=false
