build_geth_image:
	docker build -t geth ./geth_node

build_chainlink_image:
	# Pull docker images required by chainlink and then build chainlink from scratch
	docker build -t smartcontract/chainlink ../../

copy_keys:
	mkdir -p geth_node/root/keystore
	mkdir -p chainlink_node/clroot/keys
	cp ../../tools/keys/0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f.json geth_node/root/keystore/
	cp ../../tools/keys/0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f.json chainlink_node/clroot/keys

start_network:
	@echo "\nStarting Network. See logs/{chain_[12],geth}.log for logs"
	( docker-compose --compatibility up > /dev/null ) &

tear_down:
	docker-compose --compatibility down
	# Remove CL database
	sudo rm -rf chainlink_node

initial_setup: tear_down build_geth_image build_chainlink_image copy_keys start_network create_job

ETH_LOG_JOB = '$(shell cat eth_log_job.json)'
create_job:
	# wait for CL to start
	bash scripts/search_chainlink_logs.sh 'Listening and serving HTTP on port 6688'
	# add ethlog job
	bash scripts/log_in_to_chainlink_node.sh
	docker exec -it forks_chainlink chainlink jobs create ${ETH_LOG_JOB}

# Generates an ethereum transaction which creates a simple contract, which emits a log in its constructor.
create_curl_script:
	python3 scripts/construct_cmd.py > scripts/curl_cmd.sh

# Runs the first chain, where the ethlog contract is deployed
run_chain_1:
	@echo "\nSTARTING CHAIN 1"
	# make sure chainlink has actually started receiving blocks from geth
	bash scripts/search_chainlink_logs.sh 'Received new head'
	# broadcast contract creation transaction
	bash scripts/curl_cmd.sh
	# wait for chainlink to get notified about transaction
	bash scripts/search_chainlink_logs.sh 'New run triggered by ethlog'
	bash scripts/search_chainlink_logs.sh 'Run cannot continue because it lacks sufficient confirmations'
	# save log
	mkdir -p logs
	docker-compose --compatibility logs chainlink > logs/chain_1.log
	docker-compose --compatibility logs geth > logs/geth.log
	# tear down network before sufficient confirmations can be reached
	docker-compose --compatibility down

# Runs the 2nd chain, which eventually out-grows the first chain and uncles the ethlog contract
run_chain_2:
	@echo "\nSTARTING CHAIN 2"
	# 2nd chain should be younger than first, and so chainlink won't immediately save new heads
	bash scripts/search_chainlink_logs.sh 'Cannot save new head confirmation'
	# when 2nd chain gets longer, chainlink resumes saving heads
	bash scripts/search_chainlink_logs.sh 'New head resuming run'
	# will wait for head # to be 10 more than block # with contract creation
	bash scripts/search_chainlink_logs.sh 'Run cannot continue because it lacks sufficient confirmations services'
	# should eventually abort running running job
	bash scripts/search_chainlink_logs.sh 'presumably has been uncled'
	# save log
	mkdir -p logs
	docker-compose --compatibility logs chainlink > logs/chain_2.log
	docker-compose --compatibility logs geth >> logs/geth.log
	# tear down
	docker-compose --compatibility down

