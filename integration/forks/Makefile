# initial_setup: tear_down build_images start_network create_job
initial_setup: tear_down start_network create_job

build_images:
	docker build -t geth ./geth
	docker build -t smartcontract/chainlink ../../

start_network:
	@echo "\nStarting Network. See logs/{chain_[12],geth}.log for logs"
	( docker-compose up > /dev/null ) &

tear_down:
	docker-compose down
	# Remove CL database
	sudo rm -rf tmp

ETH_LOG_JOB = '$(shell cat eth_log_job.json)'
create_job:
	# wait for CL to start
	bash scripts/search_chainlink_logs 'Listening and serving HTTP on port 6688'
	# add ethlog job
	bash scripts/log_in_to_chainlink_node
	docker exec -it forks_chainlink chainlink jobs create ${ETH_LOG_JOB}

