#!/usr/bin/env bash

max_retries=${1:-5} # Default to 5 if not provided
count=0

until go run main.go download all \
  --output-dir ../ \
  --gh-token-env-var-name GITHUB_API_TOKEN \
  --cre-cli-version v0.2.0 \
  --capability-names cron \
  --capability-version v1.0.2-alpha
do
  ((count++))
  if (( count >= max_retries )); then
    echo "❌ Failed after $max_retries attempts." >&2
    exit 1
  fi
  echo "🔁 Retrying ($count/$max_retries)..." >&2
  sleep 1
done

echo "✅ Download succeeded." >&2