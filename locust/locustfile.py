import random
import time

from locust import HttpUser, task, between


class OrderUser(HttpUser):
    """負荷テスト用ユーザー"""
    wait_time = between(0.05, 0.2)
    host = "http://localhost:8080"

    @task(3)
    def get_order(self):
        """GET /api/orders/{id} - 読み取り負荷"""
        order_id = random.randint(1, 100)
        with self.client.get(
            f"/api/orders/{order_id}",
            name="/api/orders/{id}",
            catch_response=True,
        ) as response:
            if response.status_code != 200:
                response.failure(f"Status {response.status_code}")

    @task(2)
    def create_order(self):
        """POST /api/orders - 書き込み負荷（ステートレス）"""
        payload = {
            "product_name": f"負荷テスト商品_{int(time.time() * 1000)}",
            "quantity": random.randint(1, 10),
            "note": "負荷テストによる注文",
        }
        with self.client.post(
            "/api/orders",
            json=payload,
            name="/api/orders",
            catch_response=True,
        ) as response:
            if response.status_code != 201:
                response.failure(f"Status {response.status_code}")

    @task(2)
    def update_order(self):
        """PUT /api/orders/{id} - 更新負荷"""
        order_id = random.randint(1, 100)
        payload = {
            "quantity": random.randint(1, 20),
            "note": f"負荷テスト更新_{int(time.time() * 1000)}",
        }
        with self.client.put(
            f"/api/orders/{order_id}",
            json=payload,
            name="/api/orders/{id}",
            catch_response=True,
        ) as response:
            if response.status_code != 200:
                response.failure(f"Status {response.status_code}")

    @task(1)
    def confirm_order(self):
        """POST /api/orders/{id}/confirm - ステートフル負荷"""
        # Step 1: 注文を作成
        create_payload = {
            "product_name": f"確定テスト商品_{int(time.time() * 1000)}",
            "quantity": 1,
        }
        create_res = self.client.post(
            "/api/orders",
            json=create_payload,
            name="/api/orders (for confirm)",
        )
        if create_res.status_code != 201:
            return

        order_id = create_res.json()["id"]

        # Step 2: トークンを取得
        get_res = self.client.get(
            f"/api/orders/{order_id}",
            name="/api/orders/{id} (for confirm)",
        )
        if get_res.status_code != 200:
            return

        token = get_res.json()["confirmation_token"]

        # Step 3: 注文を確定
        confirm_payload = {"confirmation_token": token}
        with self.client.post(
            f"/api/orders/{order_id}/confirm",
            json=confirm_payload,
            name="/api/orders/{id}/confirm",
            catch_response=True,
        ) as response:
            if response.status_code != 200:
                response.failure(f"Status {response.status_code}")
