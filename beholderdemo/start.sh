#!/bin/sh

REPOS=~/repos
CHAINLINK=$REPOS/chainlink
BEHOLDER_DEMO=$REPOS/atlas/beholder

echo Builing chainlink-dev
make chainlink-dev || exit 1

start_beholder_stack() {
	echo "\n\nTo start consuming custom messages emitted from chainlink node to OTel Collector
	Run this command in separate terminal:

	cd $BEHOLDER_DEMO
	make consume-topic
	"

	cd $BEHOLDER_DEMO
	echo "\n\nSrating Beholder stack"
	make start
	echo "Stop Beholder Demo App (stop emitting messages to OTel Collector)"
	docker compose stop beholderdemo
	open http://localhost:3000/d/ads286ty3ah34f/beholder-demo
}

start_postgres() {
	echo "\n\nStarting Postgres"
	docker run --rm --name cl-postgres -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_USER=chainlink -e POSTGRES_DB=chainlink_unit_test -p 5433:5432 -d postgres:15.5
	echo "Check Postgress connection"
	until docker run --rm --name psql --link cl-postgres postgres:15.5 psql "postgresql://chainlink:mysecretpassword@cl-postgres:5432/chainlink_unit_test?sslmode=disable"
	do
  		echo "\nTrying to connect to Postgres"
  		sleep 1
	done
}

start_chainlink() {
	cd $CHAINLINK
	echo "\n\nSratring chainlink node"
	OTEL_SERVICE_NAME=beholderdemo ./chainlink node -c ~/.chainlink-sepolia/config.toml -s ~/.chainlink-sepolia/secrets.toml start || exit 1
}

echo "\n\nRemoving all running containers"
docker rm $(docker stop $(docker ps -aq))

start_beholder_stack

start_postgres

start_chainlink
 


