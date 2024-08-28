#!/bin/bash

REPOS=~/repos
CHAINLINK=$REPOS/chainlink
BEHOLDER_DEMO=$REPOS/atlas/beholder

echo Builing chainlink-dev
make chainlink-dev || exit 1

echo "To start consuming custom messages emitted from chainlink node to OTel Collector
Run this command in separate terminal:

  cd $BEHOLDER_DEMO
  make consume-topic
"

start_behooder_stack() {
	cd $BEHOLDER_DEMO
	echo "\n\nSrating Beholder stack"
	make start
	echo "Stop Beholder Demo App (stop emitting messages to OTel Collector)"
	docker compose stop beholderdemo
	open http://localhost:3000/d/ads286ty3ah34f/beholder-demo
}

start_chainlink() {
	cd $CHAINLINK
	echo "\n\nSratring chainlink node"
	OTEL_SERVICE_NAME=beholderdemo ./chainlink node -c ~/.chainlink-sepolia/config.toml -s ~/.chainlink-sepolia/secrets.toml start || exit 1
}

start_behooder_stack

start_chainlink
 


