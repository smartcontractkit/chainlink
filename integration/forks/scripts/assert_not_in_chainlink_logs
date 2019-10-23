num_complete=`docker-compose --compatibility logs chainlink | grep -c "$1"`;
if [ $num_complete != "0" ]; then
  echo "Found $1 in chainlink logs and shouldn't have";
  exit 1;
fi
