#!/bin/bash
# This script copies 5 jobspecs, login credentials, and a script to each of the 5 docker containers
# It will then execute the script that logins, deletes a job, and creates a new job with the new jobspecs for each container
# 
# To run it:
# ./newJob.sh (jobID to delete)
# 
docker cp jobspec_1.toml chainlink-node-1:/root
docker cp jobspec_2.toml chainlink-node-2:/root
docker cp jobspec_3.toml chainlink-node-3:/root
docker cp jobspec_4.toml chainlink-node-4:/root
docker cp jobspec_5.toml chainlink-node-5:/root
docker cp apicredentials chainlink-node-1:/root
docker cp apicredentials chainlink-node-2:/root
docker cp apicredentials chainlink-node-3:/root
docker cp apicredentials chainlink-node-4:/root
docker cp apicredentials chainlink-node-5:/root
docker cp linker.sh chainlink-node-1:/root
docker cp linker.sh chainlink-node-2:/root
docker cp linker.sh chainlink-node-3:/root
docker cp linker.sh chainlink-node-4:/root
docker cp linker.sh chainlink-node-5:/root
docker exec -it chainlink-node-1 ./linker.sh $1 jobspec_1.toml 1
docker exec -it chainlink-node-2 ./linker.sh $1 jobspec_2.toml 2
docker exec -it chainlink-node-3 ./linker.sh $1 jobspec_3.toml 3
docker exec -it chainlink-node-4 ./linker.sh $1 jobspec_4.toml 4
docker exec -it chainlink-node-5 ./linker.sh $1 jobspec_5.toml 5