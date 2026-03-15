# sysbench MySQL 負荷試験

## 概要

`sysbench` で MySQL に直接負荷をかけるサンプル。
デフォルトでは `oltp_read_write` ワークロードを実行し、結果を `results/*.log` に保存する。

## 前提条件

ルートディレクトリでDB直接負荷試験向け基盤（MySQL + 監視、APIなし）が起動済みであること（推奨）:

```bash
make -C .. sysbench-up
```

## 実行方法

```bash
# ルートから実行する場合（推奨）
make -C .. sysbench

# 1) テストデータ作成 + 2) 負荷試験実行
./run.sh all

# テストデータ作成のみ
./run.sh prepare

# 負荷試験のみ
./run.sh run

# テストデータ削除
./run.sh cleanup

# 終了
make -C .. sysbench-down
```

`make` でモードを切り替える場合:

```bash
make -C .. sysbench-prepare
make -C .. sysbench-run ARGS='--rand-type=uniform'
make -C .. sysbench-cleanup
```

追加の `sysbench` オプションは第2引数以降で渡せる:

```bash
./run.sh run --rand-type=uniform
```

## 主な環境変数

| 変数 | デフォルト | 説明 |
|------|------------|------|
| `MYSQL_HOST` | `mysql` | 接続先MySQLホスト |
| `MYSQL_PORT` | `3306` | 接続先ポート |
| `MYSQL_USER` | `app` | DBユーザー |
| `MYSQL_PASSWORD` | `password` | DBパスワード |
| `MYSQL_DB` | `orders_db` | データベース名 |
| `TEST_NAME` | `oltp_read_write` | ワークロード名（例: `oltp_read_only`） |
| `TABLES` | `4` | テーブル数 |
| `TABLE_SIZE` | `10000` | 1テーブルあたり行数 |
| `THREADS` | `16` | 実行スレッド数 |
| `TIME_SECONDS` | `60` | 実行秒数 |
| `REPORT_INTERVAL` | `2` | 中間レポート間隔（秒） |

## 例: 読み取り中心のシナリオ

```bash
TEST_NAME=oltp_read_only THREADS=32 TIME_SECONDS=120 ./run.sh all
```

## ファイル構成

```
sysbench/
├── Dockerfile          # sysbench 実行イメージ
├── docker-compose.yml  # sysbench 実行用 Compose
├── run.sh              # prepare/run/cleanup ラッパー
├── results/            # 実行ログ（run時に作成）
└── README.md
```
