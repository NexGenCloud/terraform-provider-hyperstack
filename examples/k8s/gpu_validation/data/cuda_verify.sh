#!/bin/bash

# Check the Pod's logs
LOGS=$(kubectl -n "${NAMESPACE}" logs "pod/${NAME}")
CLUSTER="${CLUSTER_NAME:-undefined}"
printf "Cluster name: %s\n" "$CLUSTER"
# Check for the success message in the logs
if echo "$LOGS" | grep -q "Test PASSED"; then
  echo "✅ GPU test pod completed successfully."
  exit 0
else
  echo "❌ GPU test pod failed or did not produce expected output."
  echo "Pod Logs:"
  echo "$LOGS"
  exit 1
fi







