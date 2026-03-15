# Artillery 負荷試験

## 概要

Artillery は YAML + JavaScript で負荷テストシナリオを記述するOSSツール。
設定ベースのアプローチで、シンプルな記述で負荷テストを定義できる。

## 可視化

- **HTMLレポート生成**: `artillery report` コマンドでJSONをHTMLに変換
- テスト完了後に `results/report.html` として出力

## 前提条件

ルートディレクトリでAPI負荷試験向け基盤が起動済みであること（推奨）:

```bash
make -C .. artillery-up
```

## 実行方法

```bash
# ルートから実行する場合（推奨）
make -C .. artillery

# 負荷試験を実行 + HTMLレポートを生成
docker compose run --rm artillery

# 終了
make -C .. artillery-down
```

## 結果の確認

```bash
# HTMLレポートをブラウザで確認
xdg-open results/report.html
```

## ファイル構成

```
artillery/
├── docker-compose.yml    # Artillery実行用Docker Compose
├── load-test.yml         # 負荷テスト設定（YAML）
├── helpers.js            # カスタムヘルパー関数（JavaScript）
├── results/
│   ├── report.json       # 実行後に生成されるJSONレポート
│   └── report.html       # 実行後に生成されるHTMLレポート
└── README.md
```
