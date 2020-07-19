#! /bin/bash

set -eux

curl https://www.googleapis.com/sql/v1beta4/projects/tag-mng-243823/instances/tag-mng-cloud/export \
	-H "Authorization: Bearer $(gcloud auth print-access-token)" \
	-H "Content-Type: application/json; charset=utf-8" \
	-d '{"exportContext":{"fileType": "SQL","uri": "gs://tag-mng/cloud_sql_dump.sql","databases": ["tag-mng"]}}'