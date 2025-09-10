#!/bin/zsh
set -euo pipefail

REGION="${REGION:-us-east-1}"

echo "Creating Users table..."
aws dynamodb create-table \
  --table-name Users \
  --attribute-definitions '[{"AttributeName":"user_id","AttributeType":"S"},{"AttributeName":"email","AttributeType":"S"}]' \
  --key-schema '[{"AttributeName":"user_id","KeyType":"HASH"}]' \
  --billing-mode PAY_PER_REQUEST \
  --global-secondary-indexes '[{"IndexName":"email-index","KeySchema":[{"AttributeName":"email","KeyType":"HASH"}],"Projection":{"ProjectionType":"ALL"}}]' \
  --region "$REGION" \
  || echo "Users table already exists or create failed."

echo "Creating ActivityLogs table..."
aws dynamodb create-table \
  --table-name ActivityLogs \
  --attribute-definitions '[{"AttributeName":"user_id","AttributeType":"S"},{"AttributeName":"timestamp","AttributeType":"N"},{"AttributeName":"activity_type","AttributeType":"S"}]' \
  --key-schema '[{"AttributeName":"user_id","KeyType":"HASH"},{"AttributeName":"timestamp","KeyType":"RANGE"}]' \
  --billing-mode PAY_PER_REQUEST \
  --global-secondary-indexes '[{"IndexName":"activity_type-index","KeySchema":[{"AttributeName":"activity_type","KeyType":"HASH"},{"AttributeName":"timestamp","KeyType":"RANGE"}],"Projection":{"ProjectionType":"ALL"}}]' \
  --region "$REGION" \
  || echo "ActivityLogs table already exists or create failed."

echo "Waiting for tables to become ACTIVE (this may take a few seconds)..."
aws dynamodb wait table-exists --table-name Users --region "$REGION"
aws dynamodb wait table-exists --table-name ActivityLogs --region "$REGION"

echo "✅ Tables created (or already existed) and are now active."
