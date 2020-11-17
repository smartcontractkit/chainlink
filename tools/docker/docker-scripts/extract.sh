grep -o '[[:P2P:]]$' keys.txt
sed -rn 's/[^,]+,([^,]+),.*/\1/p' keys.txt
perl -F, -ale 'print $F[1]' keys.txt

awk -F 'Peer ID:' '{print $2}' keys.txt
awk '/ID: [^0-9]$/ { print $2 }' keys.txt
