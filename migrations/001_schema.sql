
--create types
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'scorekind') THEN
        CREATE TYPE scorekind AS ENUM ('points', 'time');
    END IF;
    --more types here...
END$$;



CREATE TABLE IF NOT EXISTS guilds (
	guild_id BIGINT PRIMARY KEY NOT NULL,
	description VARCHAR(256) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS channels (
	channel_id BIGINT PRIMARY KEY NOT NULL,
	guild_id BIGINT NOT NULL
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	running BOOLEAN NOT NULL DEFAULT FALSE,
	UNIQUE(guild_id, channel_id)
);
CREATE INDEX IF NOT EXISTS channels_id ON channels(channel_id);

CREATE TABLE IF NOT EXISTS flags (
	flag_id SMALLINT PRIMARY KEY NOT NULL,
	abbr VARCHAR(64) NOT NULL UNIQUE,
	emoji VARCHAR(64) NOT NULL
);
CREATE INDEX IF NOT EXISTS flags_id ON flags(flag_id);
CREATE INDEX IF NOT EXISTS flags_id ON flags(abbr);

CREATE TABLE IF NOT EXISTS flag_mappings (
	guild_id BIGINT NOT NULL
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id BIGINT NOT NULL
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	flag_id SMALLINT NOT NULL
		REFERENCES flags(flag_id)
		ON DELETE CASCADE,
	emoji VARCHAR(64) NOT NULL,
	PRIMARY KEY (channel_id, flag_id)
);

CREATE TABLE IF NOT EXISTS tracking (
	id BIGSERIAL,
	message_id BIGINT PRIMARY KEY NOT NULL,
	guild_id BIGINT NOT NULL
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id BIGINT NOT NULL
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	address VARCHAR(64) NOT NULL,
	CONSTRAINT tracking_unique_address UNIQUE (guild_id, channel_id, address),
	CONSTRAINT tracking_unique_message_id UNIQUE (guild_id, channel_id, message_id)
);

CREATE TABLE IF NOT EXISTS player_count_notifications (
	guild_id BIGINT NOT NULL
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id BIGINT NOT NULL
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	message_id BIGINT NOT NULL
		REFERENCES tracking(message_id)
		ON DELETE CASCADE
		ON UPDATE CASCADE,
	user_id BIGINT NOT NULL,
	threshold SMALLINT NOT NULL
		CHECK( threshold > 0),
	PRIMARY KEY (guild_id, channel_id, message_id, user_id)
);

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


CREATE TABLE IF NOT EXISTS active_servers (
	timestamp timestamp WITH TIME ZONE NOT NULL,
	address VARCHAR(64) PRIMARY KEY NOT NULL,
	protocols jsonb NOT NULL,
	name VARCHAR(128) NOT NULL,
	gametype VARCHAR(64) NOT NULL,
	passworded BOOLEAN NOT NULL DEFAULT FALSE,
	map VARCHAR(128) NOT NULL,
	map_sha256sum char(64),
	map_size INTEGER,
	version VARCHAR(64) NOT NULL,
	max_clients SMALLINT NOT NULL,
	max_players SMALLINT NOT NULL,
	score_kind scorekind NOT NULL DEFAULT 'points'
);

CREATE TABLE IF NOT EXISTS active_server_clients (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	address VARCHAR(64) NOT NULL
		REFERENCES active_servers(address)
		ON DELETE CASCADE,
	name VARCHAR(32) NOT NULL,
	clan VARCHAR(32) NOT NULL,
	country_id SMALLINT NOT NULL
		REFERENCES flags(flag_id),
	score INTEGER NOT NULL,
	is_player BOOLEAN NOT NULL,
	team SMALLINT
);

CREATE TABLE IF NOT EXISTS prev_active_servers (
	message_id BIGINT PRIMARY KEY NOT NULL
		REFERENCES tracking(message_id)
		ON DELETE CASCADE
		ON UPDATE CASCADE,
	guild_id BIGINT NOT NULL
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id BIGINT NOT NULL
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	timestamp timestamp WITH TIME ZONE NOT NULL,
	address VARCHAR(64) NOT NULL,
	protocols jsonb NOT NULL,
	name VARCHAR(128) NOT NULL,
	gametype VARCHAR(64) NOT NULL,
	passworded BOOLEAN NOT NULL DEFAULT FALSE,
	map VARCHAR(128) NOT NULL,
	map_sha256sum char(64),
	map_size INTEGER,
	version VARCHAR(64) NOT NULL,
	max_clients SMALLINT NOT NULL,
	max_players SMALLINT NOT NULL,
	score_kind scorekind NOT NULL DEFAULT 'points'
);

CREATE TABLE IF NOT EXISTS prev_active_server_clients (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	message_id BIGINT NOT NULL
		REFERENCES tracking(message_id)
		ON DELETE CASCADE
		ON UPDATE CASCADE,
	guild_id BIGINT NOT NULL
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id BIGINT NOT NULL
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	name VARCHAR(32) NOT NULL,
	clan VARCHAR(32) NOT NULL,
	country_id SMALLINT NOT NULL
		REFERENCES flags(flag_id),
	score INTEGER NOT NULL,
	is_player BOOLEAN NOT NULL,
	team SMALLINT,
	flag_abbr VARCHAR(64) NOT NULL,
	flag_emoji VARCHAR(64) NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS prev_active_server_clients;
DROP TABLE IF EXISTS prev_active_servers;
DROP TABLE IF EXISTS active_server_clients;
DROP TABLE IF EXISTS active_servers;
DROP TABLE IF EXISTS player_count_notifications;
DROP TABLE IF EXISTS prev_message_mentions;
DROP TABLE IF EXISTS tracking;
DROP TABLE IF EXISTS flag_mappings;
DROP TABLE IF EXISTS flags;
DROP TABLE IF EXISTS channels;
DROP TABLE IF EXISTS guilds;
DROP TYPE IF EXISTS scorekind;

-- DROP TABLE IF EXISTS schema_migrations;
