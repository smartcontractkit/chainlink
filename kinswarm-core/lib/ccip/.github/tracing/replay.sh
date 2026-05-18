# Read JSON file and loop through each trace
while IFS= read -r trace; do
  curl -X POST http://localhost:3100/v1/traces \
    -H "Content-Type: application/json" \
    -d "$trace"
done < "trace-data"
