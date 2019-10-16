build_geth_image:
	docker build -t geth .


tear_down:
	docker-compose --compatibility down
	# Remove CL database
	sudo rm -rf chainlink_node/clroot
	git checkout chainlink_node/clroot/keys

start_network:
	( docker-compose --compatibility up > /dev/null ) &

initial_setup: build_geth_image start_network create_curl_script create_job

# Isolate the faster node from the network, so that it creates a blockchain
# re-org when re-attached.
isolate_geth2:
	docker network disconnect forks_gethnet forks_geth2_1 || true

ETH_LOG_JOB = '$(shell cat eth_log_job.json)'
create_job:
	sleep 3
	bash scripts/log_in_to_chainlink_node.sh
	docker exec -it forks_chainlink_1 chainlink jobs create ${ETH_LOG_JOB}

# Generates an ethereum transaction which creates a simple contract, which emits
# a log in its constructor. See "Constructing the trigger transaction" in
# `README.md`
create_curl_script:
	python3 scripts/construct_cmd.py > scripts/curl_cmd.sh

run_chain_1:
	@echo "\nSTARTING CHAIN 1"
	# make sure chainlink has actually started receiving blocks from geth
	bash scripts/search_chainlink_logs.sh 'Received new head'
	# broadcast contract creation transaction
	bash scripts/curl_cmd.sh
	# wait for chainlink to get notified about transaction
	bash scripts/search_chainlink_logs.sh 'New run triggered by ethlog'
	bash scripts/search_chainlink_logs.sh 'Run cannot continue because it lacks sufficient confirmations'
	# tear down network before sufficient confirmations can be reached
	docker-compose --compatibility down

run_chain_2:
	@echo "\nSTARTING CHAIN 2"
	# 2nd chain should be younger that first, and so chainlink won't immediately save new heads
	bash scripts/search_chainlink_logs.sh 'Cannot save new head confirmation'
	# when 2nd chain gets longer, chainlink resumes saving heads
	bash scripts/search_chainlink_logs.sh 'New head resuming run'
	# will wait for head # to be 10 more than block # with contract creation
	bash scripts/search_chainlink_logs.sh 'Run cannot continue because it lacks sufficient confirmations services'
	# should eventually abort running running job
	bash scripts/search_chainlink_logs.sh 'presumably has been uncled'
	docker-compose --compatibility down

