build_geth_image:
	docker build -t geth .

# Reset the geth and CL databases to a known-good state
clean_databases:
	sudo rm -rf ./gethnet/datadir* ../../tools/clroot
	git checkout ./gethnet/datadir*
	git checkout ../../tools/clroot
	# Clean docker-compose logs. https://github.com/docker/compose/issues/1083#issuecomment-364359851
	docker-compose --compatibility ps -q | \
		xargs -n 1 -i docker inspect --format='{{.LogPath}}' '{}' | \
		xargs -n 1 -i sudo truncate -s 0

tear_down_network:
	docker-compose --compatibility down

start_network: tear_down_network clean_databases build_geth_image
	( docker-compose --compatibility up > /dev/null ) &

# Isolate the faster node from the network, so that it creates a blockchain
# re-org when re-attached.
isolate_geth2:
	docker network disconnect forks_gethnet forks_geth2_1 || true

ETH_LOG_JOB = '$(shell cat eth_log_job.json)'
create_job:
	bash log_in_to_chainlink_node.sh
	docker exec -it forks_chainlink_1 chainlink jobs create ${ETH_LOG_JOB}

# Generates an ethereum transaction which creates a simple contract, which emits
# a log in its constructor. See "Constructing the trigger transaction" in
# `README.md`
curl_cmd.sh:
	python3 construct_cmd.py > curl_cmd.sh

# inject_transaction sends the same tx to geth1 and geth2, after geth2 has been
# isolated from the network. The goal here is to get a replay of the log emitted
# by the transaction.
inject_transaction: curl_cmd.sh create_job
	docker exec -it forks_geth_1 bash /root/curl_cmd.sh

# Keep checking until we've verified that one job was completed, and one was
# uncled.
search_for_completed_chainlink_jobs:
	./search_for_completed_chainlink_jobs.sh

