make initial_setup
make run_chain_1
make start_network
make run_chain_2

bash scripts/assert_not_in_chainlink_logs.sh 'All tasks complete for run'

echo "test passes!"
