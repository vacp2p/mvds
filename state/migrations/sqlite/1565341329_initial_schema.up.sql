CREATE TABLE mvds_states (
    type INTEGER NOT NULL,
    send_count INTEGER NOT NULL,
    send_epoch INTEGER NOT NULL,
    group_id BLOB,
    peer_id BLOB NOT NULL,
    message_id BLOB NOT NULL,
    PRIMARY KEY (message_id, peer_id)
);

CREATE INDEX idx_send_epoch ON mvds_states(send_epoch);
