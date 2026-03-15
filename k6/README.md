# Grafana k6 負荷試験

## 概要

Grafana k6 は JavaScript で負荷テストスクリプトを記述するOSSツール。
Go で実装されており、高パフォーマンスな負荷生成が可能。

## 可視化

- **組み込みWebダッシュボード**: テスト実行中にリアルタイムでメトリクスを確認
- **HTMLレポートエクスポート**: テスト完了後に `results/report.html` として出力

## 前提条件

ルートディレクトリでAPI負荷試験向け基盤が起動済みであること（推奨）:

```bash
make -C .. k6-up
```

## 実行方法

```bash
# ルートから実行する場合（推奨）
make -C .. k6

# 負荷試験を実行（完了後に results/report.html が出力される）
# --service-ports を付けることで、Webダッシュボード(5665番ポート)にホストからアクセス可能になります。
./run.sh

# 終了
make -C .. k6-down
```

## xk6-sql: SQL 直接負荷テスト

### 概要

[grafana/xk6-sql](https://github.com/grafana/xk6-sql) 拡張を使い、HTTP API を経由せず MySQL に直接 SQL クエリを発行する負荷テストサンプル。
通常の HTTP 負荷テスト (`load-test.js`) と比較することで、API レイヤーとDB レイヤーそれぞれのボトルネックを切り分けられる。

### 仕組み

- `Dockerfile.xk6` で `xk6` を使い、`xk6-sql` + `xk6-sql-driver-mysql` を組み込んだカスタム k6 バイナリをビルド
- テストスクリプト (`sql-test.js`) は VU ごとに DB コネクションを確立し、以下 2 種のクエリを実行:
  1. **Simple SELECT** (`SELECT 1`) — 接続・往復レイテンシの計測
  2. **Table SELECT** (`SELECT * FROM orders LIMIT 10`) — 実テーブルへのクエリ性能計測
- カスタムメトリクス `sql_query_duration` (Trend) / `sql_query_error_rate` (Rate) で結果を可視化

### 実行方法

```bash
# ルートから実行する場合（推奨）
make -C .. k6-sql

# DB直接負荷試験向け基盤（MySQL 等）が起動済みであること
make -C .. k6-sql-up
./run-sql.sh

# 終了
make -C .. k6-sql-down
```

> Web ダッシュボードは `http://localhost:5666` でアクセス可能（HTTP テストのポート 5665 と共存）。

### 結果の確認

```bash
# HTMLレポートをブラウザで確認
open results/report-sql.html
# または
xdg-open results/report-sql.html
```

## 結果の確認（HTTP テスト）

```bash
# HTMLレポートをブラウザで確認
open results/report.html
# または
xdg-open results/report.html
```

## ファイル構成

```
k6/
├── docker-compose.yml    # k6 実行用 Docker Compose（HTTP テスト / SQL テスト）
├── Dockerfile.xk6        # xk6-sql 拡張入りカスタム k6 ビルド
├── run.sh                # HTTP 負荷テスト実行スクリプト
├── run-sql.sh            # SQL 直接負荷テスト実行スクリプト
├── scripts/
│   ├── load-test.js      # HTTP API 負荷テストスクリプト
│   └── sql-test.js       # xk6-sql 直接 SQL 負荷テストスクリプト
├── results/
│   ├── report.html       # HTTP テストの HTML レポート
│   └── report-sql.html   # SQL テストの HTML レポート
└── README.md
```
