#!/usr/bin/env bash
#
# see https://github.com/tronprotocol/java-tron/blob/develop/docker for reference

dir="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

container_image="tronprotocol/java-tron:GreatVoyage-v4.7.4"

node_count=1

if [ -n "${CUSTOM_IMAGE:-}" ]; then
  container_image="${CUSTOM_IMAGE}"
fi

echo "Using container image: ${container_image}"

if [ $# -ne 1 ]; then
  genesis_address="TDRVFH1KLFhAmYvrXdk1hbuNQqgkVtdBX5"
else
  genesis_address="$1"
fi

set -e pipefail

bash "${dir}/java-tron.down.sh"

listen_ips=""
if [ "$(uname)" = "Darwin" ]; then
	echo "Listening on all interfaces on MacOS"
	listen_ips="0.0.0.0"
else
	docker_ip=$(docker network inspect bridge -f '{{range .IPAM.Config}}{{.Gateway}}{{end}}')
	if [ -z "${docker_ip}" ]; then
		echo "Could not fetch docker ip."
		exit 1
	fi
	listen_ips="127.0.0.1 ${docker_ip}"
fi

network_name="chainlink-tron.network"

if ! docker network inspect "$network_name" >/dev/null 2>&1; then
    docker network create --subnet=172.255.0.0/24 "$network_name"
    echo "Docker network '$network_name' created successfully."
fi

echo "Starting java-tron nodes"
echo "Genesis test account address: ${genesis_address}"

temp_dir=$(mktemp -d)

for ((i=1; i<=$node_count; i++)); do
  container_name="chainlink-tron.java-tron.$i"
  container_ip="172.255.0.10$i"

  echo "Starting ${container_name} (${container_ip})"

  if [ $i -eq 1 ]; then
    need_sync_check="false"
    startup_args="--witness"
  else
    need_sync_check="true"
    startup_args=""
  fi

  temp_conf="${temp_dir}/java-tron-$i.conf"
  sed "s/#genesis_address#/${genesis_address}/g; s/#container_ip#/${container_ip}/g; s/#need_sync_check#/${need_sync_check}/" "${dir}/java-tron.conf" > "${temp_conf}"
  echo "Created temp config: ${temp_conf}"

  full_node_http_port="${i}6667"
  solidity_node_http_port="${i}6668"
  full_node_grpc_port="${i}6669"
  solidity_node_grpc_port="${i}6670"

  listen_args=()
  for ip in $listen_ips; do
    listen_args+=("-p" "${ip}:${full_node_http_port}:16667")
    listen_args+=("-p" "${ip}:${solidity_node_http_port}:16668")
    listen_args+=("-p" "${ip}:${full_node_grpc_port}:16669")
    listen_args+=("-p" "${ip}:${solidity_node_grpc_port}:16670")
    listen_args+=("-p" "${ip}:16671:16671")
    listen_args+=("-p" "${ip}:16672:16672")
  done

  docker run \
    "${listen_args[@]}" \
    -d \
    --name "${container_name}" \
    --ip "${container_ip}" \
    --network "${network_name}" \
    --mount "type=bind,source=${temp_conf},target=/java-tron.conf" \
    --entrypoint bash \
    "${container_image}" \
    "-c" \
    "./bin/FullNode -c /java-tron.conf --debug $startup_args & mkdir -p logs && touch ./logs/tron.log && tail -F ./logs/tron.log" \

  echo "Waiting for ${container_name} container to become ready.."
  start_time=$(date +%s)
  prev_output=""
  while true; do
    output=$(docker logs "${container_name}" 2>&1)
    if [[ "${output}" != "${prev_output}" ]]; then
      echo -n "${output#$prev_output}"
      prev_output="${output}"
    fi

    if [[ $output == *"Update solid block number to 2"* ]]; then
      echo ""
      echo "${container_name} is ready."
      break

    fi

    current_time=$(date +%s)
    elapsed_time=$((current_time - start_time))

    if ((elapsed_time > 600)); then
      echo "Error: Command did not become ready within 600 seconds"
      exit 1
    fi

    sleep 3
  done
done
