CHANGED="$1"
IFS=' ' read -r -a CHANGED_ARRAY <<< "$CHANGED"
slither_errors=0
for file in "${CHANGED_ARRAY[@]}"; do
  # Select the correct solc version for the file as it might differ between contracts leading to compilation errors

#            slither --fail-none --checklist --markdown-root ${{ github.server_url }}/${{ github.repository }}/blob/${{ github.sha }}/ "$file" >> slither-results.md
    source ./contracts/scripts/select_solc_version.sh "$file"
    exit_code=$?
    if [ $exit_code -ne 0 ]; then
      echo "Error: Failed to select the correct solc version for $file."
      slither_errors=$((slither_errors+1))
      continue
    fi
      slither --fail-none --checklist "$file" >> slither-results.md  || slither_errors=$((slither_errors+1))

done

if [ $slither_errors -ne 0 ]; then
  echo
  echo "There were $slither_errors Slither errors."
  exit 1
else
  echo "All files are properly formatted."
fi
