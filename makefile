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
	docker run -d --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -v /mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/db/create_db:/docker-entrypoint-initdb.d/ postgres:11-alpine
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

# Cloud SQLからデータを取得
genclouddata:
	# GCSバケットへcurlでエクスポート(https://cloud.google.com/sql/docs/postgres/import-export/exporting?hl=ja#rest)
	scripts/sql_to_gcs.sh
	# GCSバケットからローカルへDL(https://cloud.google.com/storage/docs/downloading-objects?hl=ja#gsutil)
	sleep 3 && gsutil cp gs://tag-mng/cloud_sql_dump.sql /mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/_data/cloud_sql_dump.sql

## NOTE: apiに変更があった場合は、make buildapiでイメージを更新&GCRへpushする。で、cloud runをdestroy -> applyする
buildapi:
	docker build -t gcr.io/tag-mng-243823/api --no-cache=true -f deployments/dockerfile/api/Dockerfile .
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

# DON'T EXECUTE 'make dev_app'
# NOTE: /mnt/c配下でnpm run serveがとても遅くてつらいため、開発は~/work/tag-mng/appでやる。
dev_app:
	# まず、/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/app を~/work/tag-mng/appへ同期
	rsync -av app/src/ ~/work/tag-mng/app/src/
	# cloudsqlをローカルでプロキシして、apiを起動する
	make proxy_cloudpg
	DBDSN="host=localhost dbname=tag-mng user=xxxxx password=xxxxx sslmode=disable port=15432" SESSION_KEY="xxxxxx" DEBUG=1 go run cmd/api/main.go
	# ~/work/tag-mng/app で開発
	# 次に、~/work/tag-mng/app　を/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/app　に同期
	rsync -av ~/work/tag-mng/app/src/ app/src/
	# 最後に、/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/app で動作確認
	cd app && npm run serve

# DON'T EXECUTE 'make dev_front'
# NOTE: /mnt/c配下でnpm run serveがとても遅くてつらいため、開発は~/work/tag-mng/frontでやる。
# TODO: そもそも、/mnt/c配下のリポジトリをdebianのhomeディレクトリ配下に移動したい
dev_front:
	# まず、/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/front を~/work/tag-mng/frontへ同期
	# rsync -av front/ ~/work/tag-mng/front/
	rsync -av front/src ~/work/tag-mng/front/src
	# cloudsqlをローカルでプロキシして、apiを起動する
	make proxy_cloudpg
	DBDSN="host=localhost dbname=tag-mng user=xxxxx password=xxxxx sslmode=disable port=15432" SESSION_KEY="xxxxxx" DEBUG=1 go run cmd/api/main.go
	# ~/work/tag-mng/front で開発
	cd ~/work/tag-mng/front
	# 次に、~/work/tag-mng/front　を/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/front　に同期
	cd /mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng
	rsync -av ~/work/tag-mng/front/src/ front/src/
	# 最後に、/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/front で動作確認
	cd front && npm run serve