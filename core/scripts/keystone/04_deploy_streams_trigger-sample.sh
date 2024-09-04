#!/bin/bash
# This is for arb sepolia, see jobspec on https://cl-df-mercury-arb-sepolia-0.main.stage.cldev.sh/jobs/34/definition
go run main.go \
 deploy-streams-trigger \
 --verifierproxycontractaddress=$VERIFIER_PROXY_CONTRACT_ADDRESS \
 --verifiercontractaddress=$VERIFIER_CONTRACT_ADDRESS \
 --chainid=$CHAIN_ID \
 --fromblock=$FROM_BLOCK \
 --linkfeedid=$LINK_FEED_ID \
 --nativefeedid=$NATIVE_FEED_ID \
 --feedid=$FEED_ID \
 --dryrun=true
