#!/usr/bin/env bash

# Parse command line arguments
max_retries=${1:-5} # Default to 5 if not provided
capability_names=${2:-cron} # Default to cron if not provided
capability_version=${3:-v1.0.2-alpha} # Default to v1.0.2-alpha if not provided
output_dir=${4:-../} # Default to ../ if not provided

# Display usage if help is requested
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
  echo "Usage: $0 [max_retries] [cre_cli_version] [capability_names] [capability_version] [output_dir]"
  echo "  max_retries: Maximum number of retry attempts (default: 5)"
  echo "  capability_names: Capability names to download (default: cron)"
  echo "  capability_version: Capability version to download (default: v1.0.2-alpha)"
  echo "  output_dir: Directory to save the binaries (default: ../)"
  echo ""
  echo "Example: $0 5 v0.2.1 \"cron,automation\" v1.0.3-alpha ./binaries"
  exit 0
fi

echo "🔧 Using configuration:"
echo "  Max retries: $max_retries"
echo "  Capability names: $capability_names"
echo "  Capability version: $capability_version"
echo "  Output directory: $output_dir"
echo ""

count=0

# Build capability-names flags array
capability_flags=()
IFS=',' read -ra CAPABILITY_ARRAY <<< "$capability_names"
for capability in "${CAPABILITY_ARRAY[@]}"; do
  # Trim whitespace and remove quotes
  capability="${capability#"${capability%%[![:space:]]*}"}"  # Remove leading whitespace
  capability="${capability%"${capability##*[![:space:]]}"}"  # Remove trailing whitespace
  capability="${capability#\"}"  # Remove leading quote
  capability="${capability%\"}"  # Remove trailing quote
  if [[ -n "$capability" ]]; then
    capability_flags+=(--names "$capability")
  fi
done

until go run main.go download capabilities \
  --output-dir "$output_dir" \
  --gh-token-env-var-name GITHUB_API_TOKEN \
  "${capability_flags[@]}" \
  --version "$capability_version"
do
  ((count++))
  if (( count >= max_retries )); then
    echo "❌ Failed after $max_retries attempts." >&2
    exit 1
  fi
  echo "🔁 Retrying ($count/$max_retries)..." >&2
  sleep 30 # Wait before retrying due to rate limit issues
done

echo "✅ Download succeeded." >&2
