# 負荷試験ツール比較サンドボックス

負荷試験ツール（k6、Locust、Gatling、Artillery）を同一のAPIサーバに対して実行し、各ツールの特性を比較検証するためのサンドボックス環境。

## 構成

```
performance-testing-sandbox/
├── docker-compose.yml          # 共通基盤（API + MySQL + 監視）
├── api-server/                 # Go REST APIサーバ
├── mysql/                      # MySQL初期化SQL
├── monitoring/                 # Prometheus + Grafana + MySQL Exporter
├── k6/                         # Grafana k6 負荷試験
├── locust/                     # Locust 負荷試験
├── gatling/                    # Gatling 負荷試験
└── artillery/                  # Artillery 負荷試験
```

## 前提条件

- Docker / Docker Compose が使用可能
- Ubuntu 24.04

## クイックスタート

### 1. 共通基盤の起動

```bash
docker compose up -d
```

以下のサービスが起動します：
| サービス | URL | 説明 |
|----------|-----|------|
| API Server | http://localhost:8080 | Go REST API |
| MySQL | localhost:3306 | データベース |
| Grafana | http://localhost:3000 | MySQL負荷ダッシュボード（admin/admin） |
| Prometheus | http://localhost:9090 | メトリクス収集 |

### 2. ヘルスチェック

```bash
curl http://localhost:8080/health
```

もしくは以下のシェルスクリプトですべてのAPIのヘルスチェックを実行できる。

```bash
bash ./test_api.sh
```

### 3. 負荷試験の実行

各ツールのディレクトリに移動して実行：

```bash
# k6
cd k6 && ./run.sh && cd ..

# Locust（Web UI: http://localhost:8089）
cd locust && docker compose up -d && cd ..

# Gatling
cd gatling && docker compose run --rm gatling && cd ..

# Artillery
cd artillery && docker compose run --rm artillery && cd ..
```

各ツールの詳細は対応するディレクトリの `README.md` を参照。

### 4. MySQL負荷の確認

ブラウザで http://localhost:3000 を開き、「MySQL Overview」ダッシュボードを確認。

### 5. 後片付け

```bash
# Locustを停止(Locust起動時のみ)
cd locust && docker compose down && cd ..
```
```bash
# 共通基盤を停止
docker compose down

# データを含めて完全にリセット
docker compose down -v
```

## APIエンドポイント

| メソッド | パス | 説明 |
|----------|------|------|
| `GET` | `/api/orders/{id}` | 注文詳細を取得 |
| `PUT` | `/api/orders/{id}` | 注文情報を更新 |
| `POST` | `/api/orders` | 新規注文を作成（ステートレス） |
| `POST` | `/api/orders/{id}/confirm` | 注文を確定（ステートフル） |
| `GET` | `/health` | ヘルスチェック |

## 比較観点

| 項目 | k6 | Locust | Gatling | Artillery |
|------|-----|--------|---------|-----------|
| スクリプト言語 | JavaScript | Python | Scala | YAML + JS |
| 可視化 | Webダッシュボード + HTMLレポート | Web UI（リアルタイム） | 自動生成HTMLレポート | HTMLレポート |
| ランプアップ | ✅ stages | ✅ spawn-rate | ✅ rampUsers | ✅ rampTo |
