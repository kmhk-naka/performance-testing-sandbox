# Gatling 負荷試験

## 概要

Gatling は Scala/Java で負荷テストシナリオを記述するOSSツール。
テスト完了後に詳細なHTMLレポートを自動生成する。

## 可視化

- **自動生成HTMLレポート**: テスト完了後に `results/` ディレクトリにHTMLレポートが自動出力
  - Global/Detailsタブでリクエスト別の詳細表示
  - レスポンスタイム分布、パーセンタイル推移、RPS推移のチャート
  - アクティブユーザ数推移

## 前提条件

ルートディレクトリで共通基盤が起動済みであること:

```bash
cd .. && docker compose up -d
```

## 実行方法

```bash
# 負荷試験を実行（完了後にresultsにHTMLレポートが出力される）
docker compose run --rm gatling
```

## 結果の確認

```bash
# 最新のレポートディレクトリを確認
ls -lt results/

# HTMLレポートをブラウザで確認
xdg-open results/*/index.html
```

## ファイル構成

```
gatling/
├── docker-compose.yml        # Gatling実行用Docker Compose
├── simulations/
│   └── LoadTest.scala        # 負荷テストシミュレーション（Scala）
├── results/                  # 実行後に生成されるHTMLレポート
└── README.md
```
