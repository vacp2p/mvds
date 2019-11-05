CREATE TABLE mvds_dependencies (
    msg_id BLOB PRIMARY KEY,
    dependency BLOB NOT NULL
);

CREATE INDEX idx_dependency ON mvds_dependencies(dependency);
