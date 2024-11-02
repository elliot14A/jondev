-- queries/hash_status.sql

-- name: GetHashStatus :one
select * from hash_status
limit 1;

-- name: MarkHashAsGenerated :exec
update hash_status
set is_generated = true,
    generated_at = current_timestamp
where is_generated = false;

-- name: UpdateHashLastVerified :exec
update hash_status
set last_verified_at = current_timestamp;
