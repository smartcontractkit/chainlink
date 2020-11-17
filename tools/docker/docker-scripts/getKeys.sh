#!/bin/bash
rm keys.json
docker cp apicredentials chainlink-node-1:/root
docker cp apicredentials chainlink-node-2:/root
docker cp apicredentials chainlink-node-3:/root
docker cp apicredentials chainlink-node-4:/root
docker cp apicredentials chainlink-node-5:/root
docker exec -it chainlink-node-1 chainlink admin login -f apicredentials
docker exec -it chainlink-node-2 chainlink admin login -f apicredentials
docker exec -it chainlink-node-3 chainlink admin login -f apicredentials
docker exec -it chainlink-node-4 chainlink admin login -f apicredentials
docker exec -it chainlink-node-5 chainlink admin login -f apicredentials
echo '{"1":{ "OCR":' >> keys.json
echo "🍀🍀🍀🍀 Node 1 🍀🍀🍀🍀"
docker exec -it chainlink-node-1 chainlink --json keys ocr list   >> keys.json
echo ', "P2P":' >> keys.json
docker exec -it chainlink-node-1 chainlink --json keys p2p list >> keys.json
echo ', "ETH":' >> keys.json
docker exec -it chainlink-node-1 chainlink --json keys eth list >> keys.json
echo '},"2":{ "OCR":' >> keys.json
echo "🍀🍀🍀🍀 Node 2 🍀🍀🍀🍀"
docker exec -it chainlink-node-2 chainlink --json keys ocr list   >> keys.json
echo ', "P2P":' >> keys.json
docker exec -it chainlink-node-2 chainlink --json keys p2p list >> keys.json
echo ', "ETH":' >> keys.json
docker exec -it chainlink-node-2 chainlink --json keys eth list >> keys.json
echo '},"3":{ "OCR":' >> keys.json
echo "🍀🍀🍀🍀 Node 3 🍀🍀🍀🍀"
docker exec -it chainlink-node-3 chainlink --json keys ocr list >> keys.json
echo ', "P2P":' >> keys.json
docker exec -it chainlink-node-3 chainlink --json keys p2p list >> keys.json
echo ', "ETH":' >> keys.json
docker exec -it chainlink-node-3 chainlink --json keys eth list >> keys.json
echo '},"4":{ "OCR":' >> keys.json
echo "🍀🍀🍀🍀 Node 4 🍀🍀🍀🍀"
docker exec -it chainlink-node-4 chainlink --json keys ocr list >> keys.json
echo ', "P2P":' >> keys.json
docker exec -it chainlink-node-4 chainlink --json keys p2p list >> keys.json
echo ', "ETH":' >> keys.json
docker exec -it chainlink-node-4 chainlink --json keys eth list >> keys.json
echo '},"5":{ "OCR":' >> keys.json
echo "🍀🍀🍀🍀 Node 5 🍀🍀🍀🍀"
docker exec -it chainlink-node-5 chainlink --json keys ocr list >> keys.json
echo ', "P2P":' >> keys.json
docker exec -it chainlink-node-5 chainlink --json keys p2p list >> keys.json
echo ', "ETH":' >> keys.json
docker exec -it chainlink-node-5 chainlink --json keys eth list >> keys.json
echo '}}' >> keys.json
echo "🍀🍀🍀🍀  DONE! 🍀🍀🍀🍀"
node node.js
echo "🍀🍀🍀🍀CONFIGD!🍀🍀🍀🍀"