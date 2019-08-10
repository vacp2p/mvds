CREATE TABLE mvds_messages (
    id BLOB PRIMARY KEY,
    group_id BLOB NOT NULL,
    timestamp INTEGER NOT NULL,
    body BLOB NOT NULL
);
