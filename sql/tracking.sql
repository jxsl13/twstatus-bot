

-- name: ListGuildTrackings :many
SELECT guild_id, channel_id, address, message_id
FROM tracking
ORDER BY guild_id ASC, channel_id ASC, message_id ASC;


-- Name: ListChannelTrackings :many
SELECT guild_id, channel_id, address, message_id
FROM tracking
WHERE guild_id = ?
AND channel_id = ?
ORDER BY message_id ASC;


-- name: AddTracking :exec
INSERT INTO tracking (
    guild_id,
    channel_id,
    address,
    message_id
) VALUES (?, ?, ?, ?);


-- name: RemoveTrackingByMessageId :exec
DELETE FROM tracking
WHERE guild_id = ?
AND message_id = ?;