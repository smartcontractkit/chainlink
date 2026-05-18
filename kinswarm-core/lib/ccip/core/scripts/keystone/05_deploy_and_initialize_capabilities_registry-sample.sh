#!/bin/bash

go run main.go \
 deploy-and-initialize-capabilities-registry \
 --chainid=11155111 \
 --ethurl=$ETH_URL \
 --accountkey=$ACCOUNT_KEY \
 --craddress=$CR_ADDRESS \ // 0x0d36aAC2Fd9d6d1C1F59251be6A2B337af27C52B
