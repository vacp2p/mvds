CREATE TABLE mvds_messages (
    id        BLOB PRIMARY KEY,
    group_id  BLOB    NOT NULL,
    timestamp INTEGER NOT NULL,
    body      BLOB    NOT NULL
);

CREATE INDEX idx_group_id ON mvds_messages(group_id);


CREATE TABLE mvds_parents (
    message_id BLOB NOT NULL,
    parent_id  BLOB NOT NULL,
    FOREIGN KEY (message_id) REFERENCES mvds_messages (id)
)
