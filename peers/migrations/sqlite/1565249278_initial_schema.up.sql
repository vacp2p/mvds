CREATE TABLE mvds_peers (
    group_id BLOB NOT NULL,
    peer_id BLOB NOT NULL,
    PRIMARY KEY (group_id, peer_id) ON CONFLICT REPLACE
);
