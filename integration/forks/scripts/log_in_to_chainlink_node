login() {
    docker exec -it forks_chainlink chainlink \
           admin login -f /run/secrets/apicredentials
}

# Keep trying to log in until chainlink has started its RPC service
until login
do
    sleep 0.1
done
