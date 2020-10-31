-- +migrate Up
ALTER TABLE memos ADD COLUMN is_exposed    BOOLEAN DEFAULT false;
ALTER TABLE memos ADD COLUMN exposed_at    TIMESTAMP;

DROP TRIGGER IF EXISTS update_memos_content_trigger ON memos;
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_memos_content_function() RETURNS trigger AS $$
    BEGIN
        IF OLD.content <> NEW.content THEN
            NEW.updated_at = NOW();
            NEW.notified_cnt = 0;
            RETURN NEW;
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_memos_content_trigger BEFORE UPDATE OF content ON memos
    FOR ROW EXECUTE FUNCTION update_memos_content_function();
-- +migrate StatementEnd

-- +migrate Down
ALTER TABLE memos DROP COLUMN exposed_at;
ALTER TABLE memos DROP COLUMN is_exposed;
