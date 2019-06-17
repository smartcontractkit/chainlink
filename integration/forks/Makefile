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

inject_transactions:
	curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf8d2038504a817c800830f42408080b87e6080604052348015600f57600080fd5b507fa9f7ac292bec89f3379964cf3da3a5dec83186cd18c08abe2c0603a06fb19c2960405160405180910390a160358060496000396000f3006080604052600080fd00a165627a7a72305820753aa6316e04ef5d28155af1be96e95e873ef4ab750ebf088aaa47d11950e53f0029820fb8a03c5734a70445bb8406bb386a495ad1296cceb52b49e666ec789e63ae67a61315a071cee9b8ef3573397d9c989e0b474793e5391009386922fc586b337e6c9705f3"],"id":1}' http://172.16.1.100:8545
	docker exec -it forks_geth2_1 curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf8d2038504a817c800830f42408080b87e6080604052348015600f57600080fd5b507fa9f7ac292bec89f3379964cf3da3a5dec83186cd18c08abe2c0603a06fb19c2960405160405180910390a160358060496000396000f3006080604052600080fd00a165627a7a72305820753aa6316e04ef5d28155af1be96e95e873ef4ab750ebf088aaa47d11950e53f0029820fb8a03c5734a70445bb8406bb386a495ad1296cceb52b49e666ec789e63ae67a61315a071cee9b8ef3573397d9c989e0b474793e5391009386922fc586b337e6c9705f3"],"id":1}' http://localhost:8545

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

