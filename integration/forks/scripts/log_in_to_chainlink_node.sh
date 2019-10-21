container_name=`docker ps | grep chainlink_1 | awk '{ print $NF }'`

login() {
    docker exec -it forks_chainlink chainlink \
           admin login -f /run/secrets/node_api_credentials
}

# Keep trying to log in until chainlink has started its RPC service
until login
do
    sleep 0.1
done
