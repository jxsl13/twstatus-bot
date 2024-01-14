# twstatus-bot

This is the more user friendly variant of the Teeworlds Server Status Bot which uses the DDNet HTTP master servers for fetching its data contrary to polling each individual server on its own.

## User guide

Add this bot to your Discord server: [Click here](https://discord.com/api/oauth2/authorize?client_id=628902630617513985&permissions=18685255740480&scope=bot)

### Discord setup commands
Initially you need to specify which channel you want to allow the bot to post into.
This is done by simply executing the command `/add-channel` in the channel that the bot is supposed to write the server status messages into.
Afterwards you stay in the same channel and add tracking for your Teeworlds servers like this `/add-tracking address:123.123.123.123:8301` or for ipv6 addresses you use `/add-tracking address:[fe80::9656:d028:8652:66b6]:8303`

If you want to remove tracking, you simply delete the messages that the bot created.

When you are done with your setup, you finally need to activate the channel to be updated by the bot like this `/start` in the corresponding channel.

If you want to stop the from updating server status messages for a specific channel, you can execute the `/stop` slash command in that channel.

All of these commands provide an optional parameter called `channel` which you can use to execute all of these commands in a different channel from the channel that you want to use for posting server status updates.


## Hoster guide
In case you want to host this yourself, you can use this guide to do so.

### build instructions

There are about three ways to get you up and running.
In case that you have the Go compiler toolchain installed, you can execute:
```shell
go install github.com/jxsl13/twstatus-bot@latest
```
This downloads, compiles and installs the executable into `echo "$(go env GOPATH)/bin"`.

Alternatively you can clone this repository and build the executable yourself.
Use `make build` or `go build .` for that purpose.

Lastly, you can build and deploy the docker container using `docker compose` by calling `make deploy` (`make redeploy` and `make stop`).

### Usage
Usage:
```
Environment variables:
  TWBOT_DISCORD_TOKEN         Discord App token.
  TWBOT_SUPER_ADMINS          Comma separated list of Discord User IDs that are super admins.
  TWBOT_DB_DIR                Database directory (default: ".")
  TWBOT_DB_WAL                Enable Write-Ahead-Log for SQLite (default: "false")
  TWBOT_DISCORD_GUILD_ID      Discord Bot Owner Guild ID
  TWBOT_DISCORD_CHANNEL_ID    Discord Bot Owner ChannelID for logs
  TWBOT_POLL_INTERVAL         Poll interval for DDNet's http master server (default: "16s")
  TWBOT_LEGACY_FORMAT         Use legacy message format. If disabled, rich text embeddings will be used. (default: "false")

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
  -i, --discord-channel-id string   Discord Bot Owner ChannelID for logs
  -g, --discord-guild-id string     Discord Bot Owner Guild ID
  -t, --discord-token string        Discord App token.
  -h, --help                        help for twstatus-bot
  -l, --legacy-format               Use legacy message format. If disabled, rich text embeddings will be used.
  -p, --poll-interval duration      Poll interval for DDNet's http master server (default 16s)
  -a, --super-admins string         Comma separated list of Discord User IDs that are super admins.

Use "twstatus-bot [command] --help" for more information about a command.
```

Docker usage:
Create a `.env` file in the current directory
```shell
# mandatory parameters
TWBOT_DISCORD_TOKEN="Nj..."
TWBOT_SUPER_ADMINS="134948708277026816"
TWBOT_DISCORD_GUILD_ID="1196213160911311018"
TWBOT_DISCORD_CHANNEL_ID="1196213253102129343"

# optional parameters
# format: 1h30m5s
TWBOT_POLL_INTERVAL="16s"
```

and then execute (on Linux):
```shell
make start
```

This will build and start a docker image locally.


If you want to stop the container, you can simply execute:
```shell
make stop
```