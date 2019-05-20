-- ref: https://qiita.com/s0208_wataru/items/3f56bfea05f50d6e80f8

-- シーケンスオブジェクト作成
CREATE SEQUENCE users_id_seq START 1;
CREATE SEQUENCE memos_id_seq START 1;
CREATE SEQUENCE tags_id_seq START 1;

-- 対象カラムのデフォルト値をシーケンス値に設定(insert時に、<table名>_id_seq +1)
ALTER TABLE users ALTER id SET default nextval('users_id_seq');
ALTER TABLE memos ALTER id SET default nextval('memos_id_seq');
ALTER TABLE tags ALTER id SET default nextval('tags_id_seq');

-- テーブルと連動
-- alter sequence [sequence_name] owned by [table_name].[colname];

-- 連番の開始番号を設定
SELECT setval('users_id_seq', 1, false);
SELECT setval('memos_id_seq', 1, false);
SELECT setval('tags_id_seq', 1, false);
