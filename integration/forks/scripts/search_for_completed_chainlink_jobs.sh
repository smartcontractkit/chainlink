check_count=0;
TIMEOUT=60

until ( docker-compose --compatibility logs chainlink | grep 'presumably has been uncled' > /dev/null ) ; do
  if [ $check_count -gt $TIMEOUT ]; then
    echo 'Timed out waiting for chainlink job to be uncled';
    exit 1 ;
  fi;
  check_count=$((check_count + 1))
  sleep 1;
done

num_complete=`docker-compose --compatibility logs chainlink | grep -c 'All tasks complete for run'`;
if [ $num_complete != "0" ]; then  echo "Job was run!"; exit 1; fi
