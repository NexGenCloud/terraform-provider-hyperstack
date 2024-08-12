#!/bin/bash

CLUSTER="${CLUSTER_NAME:-undefined}"
printf "Cluster name: %s\n" "$CLUSTER"
# Verify GPU Operator is Ready
echo "🔄 Verify GPU Operator Ready -- Start"

# Wait for all GPU Operator pods to be in a 'Running' or 'Succeeded' state
while [[ "$(kubectl --namespace gpu-operator --no-headers --field-selector="status.phase!=Succeeded,status.phase!=Running" get pods | wc -l)" -ne 0 ]]; do
  sleep 10
  echo "⏳ Waiting for GPU Operator to get READY..."
done

echo "✅ Verify GPU Operator Ready -- End"
exit 0