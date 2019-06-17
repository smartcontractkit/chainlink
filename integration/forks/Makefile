build_geth_image:
	docker build -t geth .

# Reset the geth and CL databases to a known-good state
clean_databases:
	sudo rm -rf ./gethnet/datadir* ../../tools/clroot
	git checkout ./gethnet/datadir*
	git checkout ../../tools/clroot

start_network: clean_databases
	docker-compose --compatibility up

# Isolate the faster node from the network, so that it creates a blockchain
# re-org when re-attached.
isolate_geth2:
	docker network disconnect forks_gethnet forks_geth2_1

create_job:
	docker exec -it forks_chainlink_1 chainlink login -f /run/secrets/node_api
	docker exec -it forks_chainlink_1 chainlink create \
		'{"initiators": [{"type": "ethlog" }],"tasks": [{ "type": "noop", "confirmations": 10 }]}'

# Re-connect geth2 to the network, so it shares its blocks with geth1, and
# creates a blockchain re-org for the chainlink node (listening to geth1) to
# deal with.
reattach_geth2:
	docker network connect forks_gethnet forks_geth2_1

# Send a simple contract-creation tx to the two nodes, once they're isolated.
# This creates a replay
send_create_tx: isolate_geth2
	sleep 0.2 # Give docker time to disconnect geth2 from network

tear_down_network:
	docker-compose --compatibility down

