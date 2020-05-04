# リモートのDBからデータのみ取得
gendata:
	ssh ochi@ddddddo.work "cd /home/pi/tag-mng/db/bk; pg_dump -U postgres tag-mng -a > data_dump.sql"
	scp ochi@ddddddo.work:/home/pi/tag-mng/db/bk/data_dump.sql /mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/_data

# ローカルにpostgresを初期化して起動
localpg:
	# postgresコンテナ起動とともにDATABASE(tag-mng)を作成します
	docker run -d --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -v /mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/db/create_db:/docker-entrypoint-initdb.d/ postgres:12-alpine
	# 起動したlocalのpostgresコンテナのtag-mngデータベースに対してマイグレートアップします
	sleep 5 && sql-migrate up -config=db/dbconfig.yml
	# 初期化されたDBに対して、初期データを投入します
	PGPASSWORD=postgres psql -h localhost -U postgres -d tag-mng -f _data/data_dump.sql 
	# memosに対しては、created_at/updated_atを以下で初期化します
	PGPASSWORD=postgres psql -h localhost -U postgres -d tag-mng -f _data/update_time.sql
	#docker run -d --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:12-alpine

conlocalpg:
	PGPASSWORD=postgres psql -h localhost -U postgres -d tag-mng

rmlocalpg:
	docker ps -a --filter name=local-postgres -q | xargs docker rm -f

test:
	go test ./internal/... -cover -coverprofile cover.out
	go tool cover -html=cover.out -o ./cover.html

# Cloud SQLへマイグレーション
## NOTE: terraform apply後、DBのappuserのpasswordをdb/dbconfig.ymlに設定すること
DB_PASSWD=$(error please input appuser db passwd)
cloudpg:
	cloud_sql_proxy -instances=tag-mng-243823:asia-northeast1:tag-mng-cloud=tcp:15432 &
	sleep 5 && sql-migrate up -config=db/dbconfig.yml -env=production
	PGPASSWORD=$(DB_PASSWD) psql -h localhost -p 15432 -U appuser -d tag-mng -f _data/data_dump.sql
	PGPASSWORD=$(DB_PASSWD) psql -h localhost -p 15432 -U appuser -d tag-mng -f _data/update_time.sql

buildapi:
	docker build -t gcr.io/tag-mng-243823/api -f deployments/dockerfile/api/Dockerfile .
	docker push gcr.io/tag-mng-243823/api

# after 'npm run build'
deployapp:
	cd app && gcloud app deploy