# Locust 負荷試験

## 概要

Locust は Python で負荷テストシナリオを記述するOSSツール。
組み込みのWeb UIによりリアルタイムでメトリクスを可視化できる。

## 可視化

- **組み込みWeb UI**: `http://localhost:8089` でリアルタイムモニタリング
  - RPS、レスポンスタイム、ユーザ数のチャート
  - エンドポイント別の詳細統計テーブル
- **CSVエクスポート**: Web UIからCSVファイルをダウンロード可能

## 前提条件

ルートディレクトリでAPI負荷試験向け基盤が起動済みであること（推奨）:

```bash
make -C .. locust-up
```

## 実行方法

```bash
# ルートから実行する場合（推奨）
make -C .. locust-up

# Locust Web UIを起動
docker compose up -d

# ブラウザで http://localhost:8089 を開く
# Web UIでユーザ数とスポーンレートを設定して開始
#   推奨設定:
#     Number of users: 50
#     Ramp up (users/second): 1

# 停止
docker compose down

# ルートから停止する場合（推奨）
make -C .. locust-down
```

## ファイル構成

```
locust/
├── docker-compose.yml    # Locust実行用Docker Compose
├── locustfile.py         # 負荷テストスクリプト（Python）
└── README.md
```
