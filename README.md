# 負荷試験ツール比較サンドボックス

負荷試験ツールを2つの軸で比較するサンドボックス環境。

- API負荷試験: `k6` / `Locust` / `Gatling` / `Artillery`
- DB直接負荷試験: `k6-sql` / `sysbench`

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
- GNU Make が使用可能
- Ubuntu 24.04

## クイックスタート

### 1. シナリオ別の基盤起動

#### API系ツール向け（k6 / Locust / Gatling / Artillery）

```bash
make api-tools-up
```

#### DB直接負荷ツール向け（k6-sql / sysbench）

```bash
make db-tools-up
```

> 直接 `docker compose` を叩く代わりに `Makefile` ターゲットを使う運用を推奨。

シナリオごとの起動サービス:
| サービス | API負荷試験 | DB直接負荷試験 | URL |
|----------|-------------|----------------|-----|
| API Server | ✅ | - | http://localhost:8080 |
| MySQL | ✅ | ✅ | localhost:3306 |
| Grafana | ✅ | ✅ | http://localhost:3000 |
| Prometheus | ✅ | ✅ | http://localhost:9090 |

### 2. ヘルスチェック（API負荷試験シナリオ）

```bash
make health-api
```

もしくは以下のシェルスクリプトですべてのAPIのヘルスチェックを実行できる。

```bash
make test-api
```

### 3. 負荷試験の実行

#### API負荷試験ツール

```bash
# k6
make k6-up
make k6
make k6-down

# Locust（Web UI: http://localhost:8089）
make locust-up
# ...UIで実行...
make locust-down

# Gatling
make gatling-up
make gatling
make gatling-down

# Artillery
make artillery-up
make artillery
make artillery-down
```

#### DB直接負荷試験ツール

```bash
# k6-sql（MySQL 直接）
make k6-sql-up
make k6-sql
make k6-sql-down

# sysbench（MySQL 直接）
make sysbench-up
make sysbench
make sysbench-down
```

各ツールの詳細は対応するディレクトリの `README.md` を参照。

### 4. MySQL負荷の確認

ブラウザで http://localhost:3000 を開き、「MySQL Overview」ダッシュボードを確認。

### 5. 後片付け

```bash
# 共通基盤を停止（現在起動中のシナリオ）
make down

# データを含めて完全にリセット
make down-v
```

### 6. ターゲット一覧の確認

```bash
make help
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
| API負荷試験 | k6 / Locust / Gatling / Artillery | APIサーバ + その背後のDB | `make api-tools-up` |
| DB直接負荷試験 | k6-sql / sysbench | MySQLへ直接SQL | `make db-tools-up` |
