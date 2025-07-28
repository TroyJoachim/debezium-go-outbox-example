-- name: CreateOutbox :one
INSERT INTO outbox (aggregate_type, aggregate_id, type, payload)
VALUES ($1, $2, $3, $4)
RETURNING id, aggregate_type, aggregate_id, type, payload, created_at;