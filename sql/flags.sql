

-- name: ListFlags :many
SELECT flag_id, abbr, emoji FROM flags ORDER BY abbr ASC;

-- name: AddFlag :exec
REPLACE INTO flags (flag_id, abbr, emoji)
VALUES (?, ?, ?);

-- name: GetFlag :many
SELECT flag_id, abbr, emoji
FROM flags
WHERE flag_id = ?
LIMIT 1;

-- name: GetFlagByAbbr :many
SELECT flag_id, abbr, emoji
FROM flags
WHERE abbr = ?
LIMIT 1;