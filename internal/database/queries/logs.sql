-- name: LogsCreate :one
INSERT INTO logs (
    request_id,
    remote_ip,
    host,
    method,
    uri,
    user_agent,
    status,
    error,
    latency,
    latency_human,
    bytes_in,
    bytes_out
) VALUES (
    :request_id,
    :remote_ip,
    :host,
    :method,
    :uri,
    :user_agent,
    :status,
    :error,
    :latency,
    :latency_human,
    :bytes_in,
    :bytes_out
) RETURNING *;

-- name: LogsGetAll :many
SELECT method, status as response, uri as path, latency_human as response_time, timestamp as created_at
FROM logs
ORDER BY timestamp DESC;

-- name: LogsGetUniqueMethods :many
SELECT DISTINCT method
FROM logs
ORDER BY method ASC;

-- name: LogsGetBasicView :many
SELECT method, status as response, uri as path, latency_human as response_time, timestamp as created_at
FROM logs
ORDER BY timestamp DESC;

-- name: LogsGetBasicViewWithOffsetLimit :many
SELECT method, status as response, uri as path, latency_human as response_time, timestamp as created_at
FROM logs
ORDER BY timestamp DESC
LIMIT sqlc.arg(begining) OFFSET sqlc.arg(limit);

-- name: LogsGetBasicViewWithOffsetLimitAdvancedOld :many
SELECT method, status as response, uri as path, latency_human as response_time, timestamp as created_at
FROM logs
WHERE timestamp >= datetime('now', sqlc.arg(timeRange))
    AND method = sqlc.arg(method)
    AND status = sqlc.arg(responseType)
ORDER BY timestamp DESC
LIMIT sqlc.arg(begining) OFFSET sqlc.arg(limit);

-- name: LogsGetBasicViewWithOffsetLimitAdvanced :many
SELECT 
    method,
    status AS response,
    uri AS path,
    latency_human AS response_time,
    timestamp AS created_at
FROM logs
WHERE timestamp >= datetime('now', sqlc.arg(timeRange))
  AND (
        sqlc.arg(method) IS NULL 
        OR sqlc.arg(method) = '' 
        OR method = sqlc.arg(method)
      )
  AND (
        sqlc.arg(responseType) IS NULL 
        OR sqlc.arg(responseType) = '' 
        OR status = sqlc.arg(responseType)
      )
ORDER BY timestamp DESC
LIMIT sqlc.arg(begining) OFFSET sqlc.arg(limit);

-- name: LogsGetMethodStats :many
SELECT 
    method,
    COUNT(*) as count,
    AVG(latency) as avg_response_time,
    MIN(latency) as min_response_time,
    MAX(latency) as max_response_time
FROM logs
GROUP BY method
ORDER BY count DESC;

-- name: LogsGetStatusStats :many
SELECT 
    status as response,
    COUNT(*) as count,
    AVG(latency) as avg_response_time
FROM logs
GROUP BY status
ORDER BY count DESC;