

-- name: ListAllTrackings :many
SELECT guild_id, channel_id, address, message_id
FROM tracking
ORDER BY guild_id ASC, channel_id ASC, message_id ASC;


-- name: ListChannelTrackings :many
SELECT guild_id, channel_id, address, message_id
FROM tracking
WHERE guild_id = $1
AND channel_id = $2
ORDER BY message_id ASC;


-- name: AddTracking :exec
INSERT INTO tracking (
    guild_id,
    channel_id,
    address,
    message_id
) VALUES ($1, $2, $3, $4);


-- name: RemoveTrackingByMessageId :exec
DELETE FROM tracking
WHERE guild_id = $1
AND message_id = $2;