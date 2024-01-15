

-- name: ListFlags :many
SELECT flag_id, abbr, emoji FROM flags ORDER BY abbr ASC;

-- name: AddFlag :exec
INSERT INTO flags (flag_id, abbr, emoji)
VALUES ($1, $2, $3)
ON CONFLICT (flag_id) DO UPDATE
SET
	abbr = $2,
    emoji = $3;

-- name: GetFlag :many
SELECT flag_id, abbr, emoji
FROM flags
WHERE flag_id = $1
LIMIT 1;

-- name: GetFlagByAbbr :many
SELECT flag_id, abbr, emoji
FROM flags
WHERE abbr = $1
LIMIT 1;