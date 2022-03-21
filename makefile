# raspberry pi内のDBからデータのみ取得
genpidata:
	ssh ochi@ddddddo.work "cd /home/pi/tag-mng/db/bk; pg_dump -U postgres tag-mng -a > data_dump.sql"
	scp ochi@ddddddo.work:/home/pi/tag-mng/db/bk/data_dump.sql /mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/_data

# ローカルにraspberry piのDBデータのpostgresを初期化して起動
localpipg:
	# postgresコンテナ起動とともにDATABASE(tag-mng)を作成します
	docker run -d --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -v /mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/db/create_db:/docker-entrypoint-initdb.d/ postgres:12-alpine
	# 起動したlocalのpostgresコンテナのtag-mngデータベースに対してマイグレートアップします
	sleep 5 && sql-migrate up -config=db/dbconfig.yml
	# 初期化されたDBに対して、初期データを投入します
	PGPASSWORD=postgres psql -h localhost -U postgres -d tag-mng -f _data/data_dump.sql
	# memosに対しては、created_at/updated_atを以下で初期化します
	PGPASSWORD=postgres psql -h localhost -U postgres -d tag-mng -f _data/update_time.sql

# ローカルにCloud SQLのDBデータのpostgresを初期化して起動
localcloudpg:
	# postgresコンテナ起動とともにDATABASE(tag-mng)を作成します
	docker run -d --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -v /home/ochi/github.com/ddddddO/memo/db/create_db:/docker-entrypoint-initdb.d/ postgres:11-alpine
	# 初期化されたDBに対して、初期データを投入します
	sleep 5 && PGPASSWORD=postgres psql -h localhost -U postgres -d tag-mng -f _data/cloud_sql_dump.sql
	# 起動したlocalのpostgresコンテナのtag-mngデータベースに対してマイグレートアップします
	sleep 5 && sql-migrate up -config=db/dbconfig.yml


conlocalpg:
	PGPASSWORD=postgres psql -h localhost -U postgres -d tag-mng

rmlocalpg:
	docker ps -a --filter name=local-postgres -q | xargs docker rm -f

test:
	go test -v ./api/... -cover -coverprofile cover.out
	go tool cover -html=cover.out -o ./cover.html

# Cloud SQLへマイグレーション
## NOTE: terraform apply後、DBのappuserのpasswordをdb/dbconfig.ymlに設定すること
#DB_PASSWD=$(error please input appuser db passwd)
cloudpg:
	cloud_sql_proxy -instances=tag-mng-243823:asia-northeast1:tag-mng-cloud=tcp:15432 &
	sleep 5 && sql-migrate up -config=db/dbconfig.yml -env=production
	#PGPASSWORD=$(DB_PASSWD) psql -h localhost -p 15432 -U appuser -d tag-mng -f _data/data_dump.sql
	#PGPASSWORD=$(DB_PASSWD) psql -h localhost -p 15432 -U appuser -d tag-mng -f _data/update_time.sql

# NOTE: 2021-04以降のデータが取得出来てない。何か制限ある？
# Cloud SQLからデータを取得
genclouddata:
	# GCSバケットへcurlでエクスポート(https://cloud.google.com/sql/docs/postgres/import-export/exporting?hl=ja#rest)
	scripts/sql_to_gcs.sh
	# GCSバケットからローカルへDL(https://cloud.google.com/storage/docs/downloading-objects?hl=ja#gsutil)
	sleep 3 && gsutil cp gs://tag-mng/cloud_sql_dump.sql /home/ochi/github.com/ddddddO/memo/_data/cloud_sql_dump.sql

## NOTE: apiに変更があった場合は、make buildapiでイメージを更新&GCRへpushする。そして、cloud runをdestroy -> applyする
## terraform destroy -target google_cloud_run_service.api -> terraform apply
buildapi:
	docker build -t gcr.io/tag-mng-243823/api --no-cache=true -f dockerfiles/api/Dockerfile .
	docker push gcr.io/tag-mng-243823/api

deployapp:
	gcloud config set project tag-mng-243823 && \
	cd app && npm run build && gcloud app deploy -q

# このタスク実行前に、make connvmでvmにログインし、sudo supervisorctl stop exposerで一旦止めること。
# デプロイ後、sudo supervisorctl start exposer の実行を忘れないこと。
deployexp:
	# コンパイルします
	go build -o _data/exposer cmd/exposer/main.go
	# configを変更します
	gcloud config set project tag-mng-243823
	# exposerをGCEへコピーします(https://cloud.google.com/compute/docs/instances/transfer-files?hl=ja#transfergcloud)
	gcloud compute scp _data/exposer hugo-generator:/home/lbfdeatq/newmemos --zone "us-central1-a"
	# ローカルのコンパイル済みのファイルを削除します
	rm _data/exposer

connvm:
	gcloud beta compute ssh --zone "us-central1-a" "hugo-generator" --project "tag-mng-243823"

prov:
	ansible-playbook --ask-vault-pass ./ansible/playbook.yml -i ./ansible/hosts

# NOTE: SQL Workbench/Jから接続して操作する。
#       redashはdocker-composeで起動したが、cloudsqlに接続出来なかった。なのでGCE上で起動するようにする。
proxy_cloudpg:
	cloud_sql_proxy -instances=tag-mng-243823:asia-northeast1:tag-mng-cloud=tcp:15432 &

# NOTE: cloud sqlからmodelsディレクトリにコード生成
#       予め、make proxy_cloudpg
# DB_PASSWD="xxxxxxxx" make xo
xo:
	xo schema postgres://appuser:$(DB_PASSWD)@localhost:15432/tag-mng?sslmode=disable
