-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    request_id TEXT,
    remote_ip TEXT,
    host TEXT,
    method TEXT,
    uri TEXT,
    user_agent TEXT,
    status INTEGER,
    error TEXT,
    latency INTEGER, 
    latency_human TEXT,
    bytes_in INTEGER,
    bytes_out INTEGER
);


CREATE INDEX idx_logs_timestamp ON logs(timestamp);
CREATE INDEX idx_logs_status ON logs(status);
CREATE INDEX idx_logs_method ON logs(method);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_logs_method;
DROP INDEX IF EXISTS idx_logs_status;
DROP INDEX IF EXISTS idx_logs_timestamp;
DROP TABLE IF EXISTS logs;
-- +goose StatementEnd