#!/usr/bin/env bash

# APIのベースURL
BASE_URL="http://localhost:8080"

echo "====================================="
echo " API Server Health & Endpoints Check"
echo "====================================="
echo ""

# 1. ヘルスチェック
echo "=== 1. Health Check (GET /health) ==="
curl -s "${BASE_URL}/health" | jq .
echo ""

# 2. 注文詳細取得 (GET)
echo "=== 2. GET Order (GET /api/orders/1) ==="
curl -s "${BASE_URL}/api/orders/1" | jq .
echo ""

# 3. 新規注文作成 (POST - ステートレス)
echo "=== 3. Create Order (POST /api/orders) ==="
# ステートレスなPOST。作成された注文のIDを取得しておく
CREATE_RES=$(curl -s -X POST "${BASE_URL}/api/orders" \
  -H "Content-Type: application/json" \
  -d '{"product_name":"テスト商品","quantity":3,"note":"手動テストによる作成"}')

echo "${CREATE_RES}" | jq .
NEW_ORDER_ID=$(echo "${CREATE_RES}" | jq -r '.id')
echo ""

# 4. 注文情報更新 (PUT)
echo "=== 4. Update Order (PUT /api/orders/1) ==="
curl -s -X PUT "${BASE_URL}/api/orders/1" \
  -H "Content-Type: application/json" \
  -d '{"quantity":99,"note":"手動テストによる更新"}' | jq .
echo ""

# 5. 注文確定 (POST - ステートフル)
echo "=== 5. Confirm Order (POST /api/orders/{id}/confirm) ==="
echo "Target Order ID: ${NEW_ORDER_ID}"
# まずはGETしてTokenを取得する
TOKEN=$(curl -s "${BASE_URL}/api/orders/${NEW_ORDER_ID}" | jq -r '.confirmation_token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
  # 取得したTokenを使って確定APIを叩く
  curl -s -X POST "${BASE_URL}/api/orders/${NEW_ORDER_ID}/confirm" \
    -H "Content-Type: application/json" \
    -d "{\"confirmation_token\":\"${TOKEN}\"}" | jq .
else
  echo "Error: Token is missing or null."
fi
echo ""

echo "====================================="
echo " Done!"
echo "====================================="
