FROM debian

WORKDIR /tag-mng

ADD . /tag-mng

# postgresまだうまくいかない
# TODO: 以下、rubyとpostgresはコンテナ(Dockerfile)分けてPodで管理する
RUN set -x && \
    apt-get update && \
    apt-get -y install ruby && \
    gem install bundler && \
    cd web && \
    bundle install && \
    apt-get -y install postgresql && \
    cd ../db && \
    psql -f 00_create_database.sql -U postgres && \
    psql -f 01_create_table.sql -U postgres -d tag-mng && \
    psql -f 02_alter_table_sequence.sql -U postgres -d tag-mng && \
    psql -f 03_alter_table_passwd.sql -U postgres -d tag-mng

EXPOSE 4567

CMD ./tag-mng/web/launch.sh

# exec
# docker run -p 8888:4567 ddddddo/tag-mng:1.0.0
