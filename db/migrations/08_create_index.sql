-- +migrate Up
create index memos_users_id_idx on memos (users_id);
create index tags_users_id_idx on tags (users_id);

-- +migrate Down
drop index memos_users_id_idx;
drop index tags_users_id_idx;
