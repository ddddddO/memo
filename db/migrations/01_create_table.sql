-- +migrate Up
CREATE TABLE users(
    id      INTEGER NOT NULL,
    name    TEXT    NOT NULL,
    passwd  CHAR(8) NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE memos(
    id       INTEGER NOT NULL,
    subject  TEXT    NOT NULL,
    content  TEXT    NOT NULL,
    users_id INTEGER REFERENCES users(id),
    PRIMARY  KEY(id)
);

CREATE TABLE tags(
    id       INTEGER NOT NULL,
    name     TEXT    NOT NULL,
    users_id INTEGER REFERENCES users(id),
    PRIMARY KEY(id)
);

CREATE TABLE memo_tag(
    memos_id INTEGER REFERENCES memos(id),
    tags_id  INTEGER REFERENCES tags(id)
);

-- +migrate Down
DROP TABLE memo_tag;
DROP TABLE tags;
DROP TABLE memos;
DROP TABLE users;

DROP SEQUENCE memos_id_seq;
DROP SEQUENCE tags_id_seq;
DROP SEQUENCE users_id_seq;
