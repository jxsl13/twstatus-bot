# twstatus-bot

This is the more user friendly variant of the Teeworlds Server Status Bot which uses the DDNet HTTP master servers for fetching its data contrary to polling each individual server on its own.

Usage:
```
Environment variables:
  TWBOT_DISCORD_TOKEN         Discord App token.
  TWBOT_DISCORD_CHANNEL_ID    Discord Channel ID.
  TWBOT_SUPER_ADMINS          Comma separated list of Discord User IDs that are super admins.
  TWBOT_DB_DIR                Database directory (default: ".")
  TWBOT_DB_WAL                Enable Write-Ahead-Log for SQLite (default: "false")
  TWBOT_DISCORD_GUILD_ID      Discord Bot Owner Guild ID
  TWBOT_POLL_INTERVAL         Poll interval for DDNet's http master server (default: "16s")

Usage:
  twstatus-bot [flags]
  twstatus-bot [command]

Available Commands:
  completion  Generate completion script
  help        Help about any command

Flags:
  -c, --config string               .env config file path (or via env variable TWBOT_CONFIG)
  -d, --db-dir string               Database directory (default ".")
  -w, --db-wal                      Enable Write-Ahead-Log for SQLite
  -i, --discord-channel-id string   Discord Channel ID.
  -g, --discord-guild-id string     Discord Bot Owner Guild ID
  -t, --discord-token string        Discord App token.
  -h, --help                        help for twstatus-bot
  -p, --poll-interval duration      Poll interval for DDNet's http master server (default 16s)
  -a, --super-admins string         Comma separated list of Discord User IDs that are super admins.
```