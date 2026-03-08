# Grafana k6 負荷試験

## 概要

Grafana k6 は JavaScript で負荷テストスクリプトを記述するOSSツール。
Go で実装されており、高パフォーマンスな負荷生成が可能。

## 可視化

- **組み込みWebダッシュボード**: テスト実行中にリアルタイムでメトリクスを確認
- **HTMLレポートエクスポート**: テスト完了後に `results/report.html` として出力

## 前提条件

ルートディレクトリで共通基盤が起動済みであること:

```bash
cd .. && docker compose up -d
```

## 実行方法

```bash
# 負荷試験を実行（完了後に results/report.html が出力される）
# --service-ports を付けることで、Webダッシュボード(5665番ポート)にホストからアクセス可能になります。
docker compose run --rm --service-ports k6
```

## 結果の確認

```bash
# HTMLレポートをブラウザで確認
open results/report.html
# または
xdg-open results/report.html
```

## ファイル構成

```
k6/
├── docker-compose.yml    # k6実行用Docker Compose
├── scripts/
│   └── load-test.js      # 負荷テストスクリプト
├── results/
│   └── report.html       # 実行後に生成されるHTMLレポート
└── README.md
```
