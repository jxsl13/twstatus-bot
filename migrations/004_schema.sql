

DROP TABLE IF EXISTS prev_message_mentions;

ALTER TABLE IF EXISTS player_count_notifications
RENAME TO player_count_notification_requests;


---- create above / drop below ----


CREATE TABLE IF NOT EXISTS prev_message_mentions (
	guild_id BIGINT NOT NULL
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id BIGINT NOT NULL
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	message_id BIGINT NOT NULL,
	user_id BIGINT NOT NULL,
	PRIMARY KEY (guild_id, channel_id, message_id, user_id)
);

ALTER TABLE IF EXISTS player_count_notification_requests
RENAME TO player_count_notifications;