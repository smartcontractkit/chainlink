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
inject_transaction: curl_cmd.sh isolate_geth2 create_job
	docker exec -it forks_geth1_1 bash /root/curl_cmd.sh
	docker exec -it forks_geth2_1 bash /root/curl_cmd.sh

# Check the logs for the geth docker containers, until geth2 has a heavier tail.
# This happens pretty fast, with overwhelming probability, because geth2 is
# running twice as fast as geth1. (See the `cpus` fields in
# docker-compose.yaml.)
ensure_geth2_has_heavier_tail:
	until python3 check_geth2_has_heavier_chain.py ; do sleep 0.5; done

# Re-connect geth2 to the network, so it shares its blocks with geth1, and
# creates a blockchain re-org for the chainlink node (listening to geth1) to
# deal with.
reattach_geth2: inject_transaction ensure_geth2_has_heavier_tail
	docker network connect forks_gethnet forks_geth2_1

# Check logs for geth1 until it's resynchronizing with geth2
await_geth1_resync: reattach_geth2
	until docker-compose --compatibility logs geth1 | grep 'Block synchronisation started'; do \
		sleep 0.5; \
	done

# Keep checking until we've verified that one job was completed, and one was
# uncled.
search_for_completed_chainlink_jobs: await_geth1_resync
	# check_count lengthens by 1 roughly each second. Used to track a timeout.
	check_count=''; \
	until ( docker-compose --compatibility logs chainlink | grep 'Finished processing task noop' ) ; do \
		if [ $${#check_count} -gt 120 ]; then \
			echo 'Timed out waiting for chainlink job to finish'; \
			exit 1 ; \
		fi; \
		check_count=$${check_count}'.'; \
		sleep 1; \
	done; \
	until ( docker-compose --compatibility logs chainlink | grep 'presumably has been uncled' ) ; do \
		if [ $${#check_count} -gt 120 ]; then \
			echo 'Timed out waiting for chainlink job to be uncled'; \
			exit 1 ; \
		fi; \
		check_count=$${check_count}'.'; \
		sleep 1; \
	done
	num_complete=`docker-compose --compatibility logs chainlink | grep -c 'All tasks complete for run'`; \
	if [ $num_complete != "1" ]; then  echo "More than one run completed!"; exit 1; fi

