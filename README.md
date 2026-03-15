# 負荷試験ツール比較サンドボックス

負荷試験ツールを2つの軸で比較するサンドボックス環境。

- API負荷試験: `k6` / `Locust` / `Gatling` / `Artillery`
- DB直接負荷試験: `sysbench`

## 構成

```
performance-testing-sandbox/
├── docker-compose.yml          # 基盤プロファイル（api / db / monitoring）
├── api-server/                 # Go REST APIサーバ
├── mysql/                      # MySQL初期化SQL
├── monitoring/                 # Prometheus + Grafana + MySQL Exporter
├── k6/                         # Grafana k6 負荷試験
├── locust/                     # Locust 負荷試験
├── gatling/                    # Gatling 負荷試験
├── artillery/                  # Artillery 負荷試験
└── sysbench/                   # sysbench MySQL負荷試験
```

## 前提条件

- Docker / Docker Compose が使用可能
- Ubuntu 24.04

## クイックスタート

### 1. シナリオ別の基盤起動

#### API負荷試験向け（API + MySQL + 監視）

```bash
docker compose --profile api --profile monitoring up -d
```

#### DB直接負荷試験向け（MySQL + 監視、APIなし）

```bash
docker compose --profile db --profile monitoring up -d
```

#### 監視基盤のみ確認したい場合（MySQL + 監視、APIなし）

```bash
docker compose --profile monitoring up -d
```

> このリポジトリのルート `docker-compose.yml` はプロファイル前提のため、`docker compose up -d` 単体ではサービスは起動しません。

シナリオごとの起動サービス:
| サービス | API負荷試験 | DB直接負荷試験 | URL |
|----------|-------------|----------------|-----|
| API Server | ✅ | - | http://localhost:8080 |
| MySQL | ✅ | ✅ | localhost:3306 |
| Grafana | ✅ | ✅ | http://localhost:3000 |
| Prometheus | ✅ | ✅ | http://localhost:9090 |

### 2. ヘルスチェック（API負荷試験シナリオ）

```bash
curl http://localhost:8080/health
```

もしくは以下のシェルスクリプトですべてのAPIのヘルスチェックを実行できる。

```bash
bash ./test_api.sh
```

### 3. 負荷試験の実行

#### API負荷試験ツール

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

#### DB直接負荷試験ツール

```bash
# sysbench（MySQL 直接負荷）
cd sysbench && ./run.sh all && cd ..
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

## 比較軸

| 軸 | ツール | 対象 | 前提基盤 |
|----|--------|------|----------|
| API負荷試験 | k6 / Locust / Gatling / Artillery | APIサーバ + その背後のDB | `--profile api`（必要に応じて `--profile monitoring`） |
| DB直接負荷試験 | sysbench | MySQLへ直接SQL | `--profile db`（必要に応じて `--profile monitoring`） |
