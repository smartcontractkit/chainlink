TIMEOUT=60

until ( docker-compose --compatibility logs chainlink | grep 'Received new head' > /dev/null ) ; do
  if [ $check_count -gt $TIMEOUT ]; then
    echo 'Timed out waiting for chainlink to receive first block';
    exit 1 ;
  fi;
  check_count=$((check_count + 1))
  sleep 1;
done
