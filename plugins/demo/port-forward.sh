#!/bin/bash
#kill ports

ps aux | grep "port-forward" | grep -v "grep" | awk '{print $2}' | while read line; do kill $line; echo "killed $line"; done
nmspc="${1:-solana-soak-67a7f}"

echo "namespace $nmspc"
kubectl get pods -n $nmspc
solPod=$(kubectl get pods -n $nmspc | grep sol | awk '{print $1}')
echo "solana pod $solPod"

#forward solana explorer
kubectl port-forward -n "$nmspc"  "$solPod" 8899:8899 &

# forward ui of node 1
kubectl port-forward -n $nmspc chainlink-0-1 6681:6688 &
echo "forwarding node 1 ui to localhost:6681"
echo "forwarding solana explorer to localhost:8899"
open -a "Brave Browser.app" "https://explorer.solana.com/?cluster=custom&customUrl=http%3A%2F%2Flocalhost%3A8899"
open -a "Brave Browser.app" "http://localhost:6681"

echo "watching node 1 pod processes"
kubectl exec --stdin --tty -n solana-soak-179f4 chainlink-0-1 --container node -- watch 'ps -aux | grep chainlink | grep -v grep' 