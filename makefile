# ローカルにpostgresをテーブル初期化して起動
localpg:
	docker run -d --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -v /mnt/c/DEV/workspace/GO/src/github.com/ddddddO/tag-mng/db:/docker-entrypoint-initdb.d/ postgres:12-alpine

rmlocalpg:
	docker ps -a --filter name=local-postgres -q | xargs docker rm -f