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

start_network:
	( docker-compose --compatibility up > /dev/null ) &

initial_setup: tear_down_network clean_databases build_geth_image
	( docker-compose --compatibility up > /dev/null ) &

# Isolate the faster node from the network, so that it creates a blockchain
# re-org when re-attached.
isolate_geth2:
	docker network disconnect forks_gethnet forks_geth2_1 || true

ETH_LOG_JOB = '$(shell cat eth_log_job.json)'
create_job:
	bash scripts/log_in_to_chainlink_node.sh
	docker exec -it forks_chainlink_1 chainlink jobs create ${ETH_LOG_JOB}

# Generates an ethereum transaction which creates a simple contract, which emits
# a log in its constructor. See "Constructing the trigger transaction" in
# `README.md`
curl_cmd.sh:
	python3 scripts/construct_cmd.py > scripts/curl_cmd.sh

# inject_transaction sends the same tx to geth1 and geth2, after geth2 has been
# isolated from the network. The goal here is to get a replay of the log emitted
# by the transaction.
inject_transaction: curl_cmd.sh create_job
	bash scripts/curl_cmd.sh

wait_for_chainlink_first_block:
	bash scripts/search_chainlink_logs.sh 'Received new head'

run_chain_1: curl_cmd.sh create_job
	bash scripts/search_chainlink_logs.sh 'Received new head'
	bash scripts/curl_cmd.sh
	bash scripts/search_chainlink_logs.sh 'New run triggered by ethlog'
	bash scripts/search_chainlink_logs.sh 'Run cannot continue because it lacks sufficient confirmations'
	docker-compose --compatibility down

run_chain_2: start_network
	bash scripts/search_chainlink_logs.sh 'Cannot save new head confirmation'
	bash scripts/search_chainlink_logs.sh 'New head resuming run'
	bash scripts/search_chainlink_logs.sh 'Run cannot continue because it lacks sufficient confirmations services'
	bash scripts/search_chainlink_logs.sh 'presumably has been uncled'
	docker-compose --compatibility down
	bash scripts/assert_not_in_chainlink_logs.sh 'All tasks complete for run'
