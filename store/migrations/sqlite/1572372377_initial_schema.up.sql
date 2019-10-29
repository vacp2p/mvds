CREATE TABLE mvds_messages (
    id        BLOB PRIMARY KEY,
    group_id  BLOB    NOT NULL,
    timestamp INTEGER NOT NULL,
    body      BLOB    NOT NULL
);

CREATE TABLE mvds_parents (
    message BLOB NOT NULL,
    parent  BLOB NOT NULL,
    FOREIGN KEY (message) REFERENCES mvds_messages (id)
)
