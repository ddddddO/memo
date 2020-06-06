-- +migrate Up
ALTER TABLE memos ADD COLUMN created_at TIMESTAMP;
ALTER TABLE memos ADD COLUMN updated_at TIMESTAMP;
ALTER TABLE memos ADD COLUMN notified_cnt INTEGER DEFAULT 0;

-- NOTE: メモ詳細の内容更新時のupdated_at/notified_cnt更新
-- ref: https://www.postgresql.jp/document/11/html/plpgsql-trigger.html
-- ref: https://github.com/rubenv/sql-migrate#writing-migrations
-- +migrate StatementBegin
CREATE FUNCTION update_memos_content_function() RETURNS trigger AS $$
    BEGIN
        IF OLD.content <> NEW.content THEN
            NEW.updated_at = NOW();
            NEW.notified_cnt = 0;
            RETURN NEW;
        END IF;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_memos_content_trigger BEFORE UPDATE OF content ON memos
    FOR ROW EXECUTE FUNCTION update_memos_content_function();
-- +migrate StatementEnd

-- NOTE: メモ詳細新規作成時のcreated_at/updated_at更新
-- +migrate StatementBegin
CREATE FUNCTION update_created_updated_at_function() RETURNS trigger AS $$
    BEGIN
        NEW.created_at = NOW();
        NEW.updated_at = NEW.created_at;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_created_updated_at_trigger BEFORE INSERT ON memos
    FOR ROW EXECUTE FUNCTION update_created_updated_at_function();
-- +migrate StatementEnd

-- +migrate Down
DROP TRIGGER IF EXISTS update_created_updated_at_trigger ON memos;
DROP FUNCTION update_created_updated_at_function();
DROP TRIGGER IF EXISTS update_memos_content_trigger ON memos;
DROP FUNCTION update_memos_content_function();
ALTER TABLE memos DROP COLUMN notified_cnt;
ALTER TABLE memos DROP COLUMN updated_at;
ALTER TABLE memos DROP COLUMN created_at;
