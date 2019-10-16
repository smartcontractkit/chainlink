check_count=0;
TIMEOUT=20

until ( docker-compose --compatibility logs chainlink | grep "$1" > /dev/null) ; do
  if [ $check_count -gt $TIMEOUT ]; then
    echo "Timed out searching chainlink logs for $1";
    exit 1 ;
  fi;
  check_count=$((check_count + 1))
  sleep 1;
done
