-- +migrate Up
DROP TRIGGER IF EXISTS update_memos_subject_trigger ON memos;

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_memos_subject_function() RETURNS trigger AS $$
    BEGIN
        IF OLD.subject <> NEW.subject THEN
            NEW.updated_at = NOW();
            NEW.notified_cnt = 0;
            RETURN NEW;
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_memos_subject_trigger BEFORE UPDATE OF subject ON memos
    FOR ROW EXECUTE FUNCTION update_memos_subject_function();
-- +migrate StatementEnd
