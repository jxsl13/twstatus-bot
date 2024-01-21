
CREATE TABLE IF NOT EXISTS player_count_notification_messages (
    channel_id BIGINT NOT NULL PRIMARY KEY
        references channels(channel_id)
        ON DELETE CASCADE,
    message_id BIGINT NOT NULL
);


---- create above / drop below ----

DROP TABLE IF EXISTS player_count_notification_messages;