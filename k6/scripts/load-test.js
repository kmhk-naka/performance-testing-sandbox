import http from 'k6/http';
import { check, sleep } from 'k6';

// ランプアップパターン
export const options = {
    stages: [
        { duration: '30s', target: 10 },  // 0 → 10 VUs
        { duration: '1m', target: 50 },  // 10 → 50 VUs
        { duration: '2m', target: 50 },  // 50 VUs 維持
        { duration: '30s', target: 0 },   // 50 → 0 VUs
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'],
        http_req_failed: ['rate<0.1'],
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// --- シナリオ: GET /api/orders/{id} ---
export function getOrder() {
    const id = Math.floor(Math.random() * 100) + 1;
    const res = http.get(`${BASE_URL}/api/orders/${id}`);
    check(res, {
        'GET /api/orders/{id} status is 200': (r) => r.status === 200,
        'GET /api/orders/{id} has id': (r) => r.status === 200 && r.json() && typeof r.json().id !== 'undefined',
    });
    sleep(0.1);
}

// --- シナリオ: POST /api/orders（ステートレス） ---
export function createOrder() {
    const payload = JSON.stringify({
        product_name: `負荷テスト商品_${Date.now()}`,
        quantity: Math.floor(Math.random() * 10) + 1,
        note: '負荷テストによる注文',
    });
    const params = { headers: { 'Content-Type': 'application/json' } };
    const res = http.post(`${BASE_URL}/api/orders`, payload, params);
    check(res, {
        'POST /api/orders status is 201': (r) => r.status === 201,
    });
    sleep(0.1);
}

// --- シナリオ: PUT /api/orders/{id} ---
export function updateOrder() {
    const id = Math.floor(Math.random() * 100) + 1;
    const payload = JSON.stringify({
        quantity: Math.floor(Math.random() * 20) + 1,
        note: `負荷テスト更新_${Date.now()}`,
    });
    const params = { headers: { 'Content-Type': 'application/json' } };
    const res = http.put(`${BASE_URL}/api/orders/${id}`, payload, params);
    check(res, {
        'PUT /api/orders/{id} status is 200': (r) => r.status === 200,
    });
    sleep(0.1);
}

// --- シナリオ: POST /api/orders/{id}/confirm（ステートフル） ---
export function confirmOrder() {
    // Step 1: 注文を作成
    const createPayload = JSON.stringify({
        product_name: `確定テスト商品_${Date.now()}`,
        quantity: 1,
    });
    const params = { headers: { 'Content-Type': 'application/json' } };
    const createRes = http.post(`${BASE_URL}/api/orders`, createPayload, params);

    if (createRes.status !== 201) {
        return;
    }
    const orderId = JSON.parse(createRes.body).id;

    // Step 2: 注文詳細を取得してトークンを取得
    const getRes = http.get(`${BASE_URL}/api/orders/${orderId}`);
    if (getRes.status !== 200) {
        return;
    }
    const token = JSON.parse(getRes.body).confirmation_token;

    // Step 3: 注文を確定
    const confirmPayload = JSON.stringify({
        confirmation_token: token,
    });
    const confirmRes = http.post(`${BASE_URL}/api/orders/${orderId}/confirm`, confirmPayload, params);
    check(confirmRes, {
        'POST /api/orders/{id}/confirm status is 200': (r) => r.status === 200,
        'POST /api/orders/{id}/confirm status is confirmed': (r) => r.status === 200 && r.json() && r.json().status === 'confirmed',
    });
    sleep(0.1);
}

// デフォルトのシナリオ: 全エンドポイントをランダムに実行
export default function () {
    const scenarios = [getOrder, createOrder, updateOrder, confirmOrder];
    const scenario = scenarios[Math.floor(Math.random() * scenarios.length)];
    scenario();
}
