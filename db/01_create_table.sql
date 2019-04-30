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
    id      INTEGER NOT NULL,
    name    TEXT    NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE memo_tag(
    memos_id INTEGER REFERENCES memos(id),
    tags_id  INTEGER REFERENCES tags(id)
);
