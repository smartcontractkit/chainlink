#!/bin/bash

go run main.go \
 deploy-contracts \
 --ocrfile=ocr_config.json \
 --chainid=1337 \
 --ethurl=ETH_URL \
 --accountkey=ACCOUNT_KEY \
 --onlysetconfig=false \
 --skipfunding=false \
 --dryrun=false
