package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/db"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
	"github.com/jxsl13/twstatus-bot/utils"
	"github.com/puzpuzpuz/xsync/v3"
)

const (
	channelOptionName = "channel"
)

var (
	reactionPlayerCountNotificationMap = map[discord.APIEmoji]int{
		discord.NewAPIEmoji(0, "1️⃣"): 1,
		discord.NewAPIEmoji(0, "2️⃣"): 2,
		discord.NewAPIEmoji(0, "3️⃣"): 3,
		discord.NewAPIEmoji(0, "4️⃣"): 4,
		discord.NewAPIEmoji(0, "5️⃣"): 5,
		discord.NewAPIEmoji(0, "6️⃣"): 6,
		discord.NewAPIEmoji(0, "7️⃣"): 7,
		discord.NewAPIEmoji(0, "8️⃣"): 8,
		discord.NewAPIEmoji(0, "9️⃣"): 9,
		discord.NewAPIEmoji(0, "🔟"):   10,
	}
	reactionPlayerCountNotificationReverseMap = map[int]discord.APIEmoji{
		1:  discord.NewAPIEmoji(0, "1️⃣"),
		2:  discord.NewAPIEmoji(0, "2️⃣"),
		3:  discord.NewAPIEmoji(0, "3️⃣"),
		4:  discord.NewAPIEmoji(0, "4️⃣"),
		5:  discord.NewAPIEmoji(0, "5️⃣"),
		6:  discord.NewAPIEmoji(0, "6️⃣"),
		7:  discord.NewAPIEmoji(0, "7️⃣"),
		8:  discord.NewAPIEmoji(0, "8️⃣"),
		9:  discord.NewAPIEmoji(0, "9️⃣"),
		10: discord.NewAPIEmoji(0, "🔟"),
	}
)

var ownerCommandList = []api.CreateCommandData{
	{
		Name:        "list-guilds",
		Description: "List all guilds that are allowed to use this bot",
	},
	{
		Name:        "add-guild",
		Description: "Allow a guild to use this bot",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "description",
				Description: "A description for this guild.",
				MinLength:   option.NewInt(4),
				MaxLength:   option.NewInt(256),
				Required:    true,
			},
			&discord.StringOption{
				OptionName:  "id",
				Description: "The guild id of the guild you want to add.",
				MinLength:   option.NewInt(1),
				MaxLength:   option.NewInt(20),
				Required:    false,
			},
		},
	},
	{
		Name:        "remove-guild",
		Description: "Remove a guild from the allowed guilds",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "id",
				Description: "The guild id of the guild you want to remove.",
				MinLength:   option.NewInt(1),
				MaxLength:   option.NewInt(20),
				Required:    false,
			},
		},
	},
	{
		Name:        "update-servers",
		Description: "Update the server list",
	},
	{
		Name:        "update-messages",
		Description: "Update all discord messages with server status lists",
	},
}

var userCommandList = []api.CreateCommandData{
	{
		Name:        "help",
		Description: "Show this help message",
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
	},
	{
		Name:           "add-channel",
		Description:    "Add a channel to the allowed channels",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to add.",
				Required:    false,
			},
		},
	},
	{
		Name:           "remove-channel",
		Description:    "Remove a channel from the allowed channels",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to remove.",
				Required:    false,
			},
		},
	},
	{
		Name:           "list-channels",
		Description:    "List all channels of the current guild that are registered for this bot",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
	},
	{
		Name:           "list-flag-mappings",
		Description:    "List all flag mappings for the current or given channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to list the flag mappings for.",
				Required:    false,
			},
		},
	},
	{
		Name:           "add-flag-mapping",
		Description:    "Add a flag mapping for the current channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "abbr",
				Description: "The abbreviation of the flag you want to add a different emoji for.",
				Required:    true,
				MinLength:   option.NewInt(2),
				MaxLength:   option.NewInt(7), // len("default")
			},
			&discord.StringOption{
				OptionName:  "emoji",
				Description: "The emoji you want to use for this flag (any text).",
				Required:    true,
				MinLength:   option.NewInt(1), // :X:
				MaxLength:   option.NewInt(256),
			},
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to add a flag mapping for.",
				Required:    false,
			},
		},
	},
	{
		Name:           "remove-flag-mapping",
		Description:    "Remove a flag mapping for the current or provided channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "abbr",
				Description: "The abbreviation of the flag you want to remove a mapping for.",
				Required:    true,
				MinLength:   option.NewInt(2),
				MaxLength:   option.NewInt(7), // len("default")
			},
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to remove a flag mapping for.",
				Required:    false,
			},
		},
	},
	{
		Name:           "list-flags",
		Description:    "show all known flags",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
	},
	{
		Name:           "add-tracking",
		Description:    "Add tracking of a Teeworlds server for the current or given channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "address",
				Description: "One or a list of comma separated server addresses that you want to track.",
				Required:    true,
				MinLength:   option.NewInt(9),
			},
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to track the server for.",
				Required:    false,
			},
		},
	},
	{
		Name:           "start",
		Description:    "Start the bot for  the given channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to start the bot for.",
				Required:    false,
			},
		},
	},
	{
		Name:           "stop",
		Description:    "Stop the bot for the given channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to stop the bot for.",
				Required:    false,
			},
		},
	},
}

type Bot struct {
	ctx             context.Context
	state           *state.State
	db              *db.DB
	superAdmins     []discord.UserID
	useEmbeds       bool
	guildID         discord.GuildID
	channelID       discord.ChannelID
	userID          discord.UserID
	c               chan model.ChangedServerStatus
	pollingInterval time.Duration
	conflictMap     *xsync.MapOf[model.MessageTarget, Backoff]
}

// New requires a discord bot token and returns a Bot instance.
// A bot token starts with Nj... and can be obtained from the discord developer portal.
func New(
	ctx context.Context,
	token string,
	db *db.DB,
	superAdmins []discord.UserID,
	guildID discord.GuildID,
	channelID discord.ChannelID,
	pollingInterval time.Duration,
	legacyMessageFormat bool,
) (*Bot, error) {

	s := state.New("Bot " + token)
	app, err := s.CurrentApplication()
	if err != nil {
		return nil, fmt.Errorf("failed to get current application: %w", err)
	}

	bot := &Bot{
		ctx:             ctx,
		state:           s,
		db:              db,
		superAdmins:     superAdmins,
		useEmbeds:       !legacyMessageFormat,
		c:               make(chan model.ChangedServerStatus, 1024),
		conflictMap:     xsync.NewMapOf[model.MessageTarget, Backoff](),
		pollingInterval: pollingInterval,
		guildID:         guildID,
		channelID:       channelID,
	}

	s.AddIntents(
		gateway.IntentGuilds | gateway.IntentGuildMessages | gateway.IntentGuildMessageReactions,
	)

	var startupOnce sync.Once
	s.AddHandler(func(*gateway.ReadyEvent) {
		// it's possible that the bot occasionally looses the gateway connection
		// calling this heavy weight function on every reconnect is not ideal
		startupOnce.Do(func() {
			me, err := s.Me()
			if err != nil {
				log.Fatalf("failed to get bot user: %v", err)
			}
			bot.userID = me.ID

			log.Println("connected to the gateway as", me.Tag())
			src, dst, err := bot.updateServers()
			if err != nil {
				log.Printf("failed to initialize server list: %v", err)
			} else {
				log.Printf("initialized server list with %d source and %d target servers", src, dst)
			}

			// sync trackings and player notification requests
			err = bot.syncDatabaseState(ctx)
			if err != nil {
				log.Fatalf("failed to synchronize database with discord state: %v", err)
			}

			// start polling
			go bot.cacheCleanup()
			go bot.serverUpdater(pollingInterval)
			for i := 0; i < max(2*runtime.NumCPU(), 5); i++ {
				go bot.messageUpdater(i + 1)
			}
		})
	})

	// requires guild message intents
	s.AddHandler(bot.handleMessageDeletion)
	s.AddHandler(bot.handleAddGuild)
	s.AddHandler(bot.handleRemoveGuild)
	s.AddHandler(bot.handleAddReactions)
	s.AddHandler(bot.handleRemoveReactions)

	r := cmdroute.NewRouter()

	// bot owner commands
	r.AddFunc("list-guilds", bot.listGuilds)
	r.AddFunc("add-guild", bot.addGuildCommand)
	r.AddFunc("remove-guild", bot.removeGuildCommand)
	r.AddFunc("update-servers", bot.updateServerListCommand)

	// user commands
	r.AddFunc("help", bot.help)
	r.AddFunc("list-channels", bot.listChannels)
	r.AddFunc("add-channel", bot.addChannel)
	r.AddFunc("remove-channel", bot.removeChannel)
	r.AddFunc("list-flags", bot.listFlags)
	r.AddFunc("add-flag-mapping", bot.addFlagMapping)
	r.AddFunc("list-flag-mappings", bot.listFlagMappings)
	r.AddFunc("remove-flag-mapping", bot.removeFlagMapping)
	r.AddFunc("add-tracking", bot.addTracking)
	r.AddFunc("start", bot.startChannel)
	r.AddFunc("stop", bot.stopChannel)

	s.AddInteractionHandler(r)

	_, err = s.BulkOverwriteGuildCommands(app.ID, bot.guildID, ownerCommandList)
	if err != nil {
		return nil, err
	}

	// update user facing commands
	err = cmdroute.OverwriteCommands(s, userCommandList)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (b *Bot) Connect(ctx context.Context) error {
	return b.state.Connect(ctx)
}

func (b *Bot) Close() error {
	return errors.Join(
		b.state.Close(),
	)
}

func (b *Bot) IsSuperAdmin(data cmdroute.CommandData) bool {
	// must be infigured guild and channel
	if !(data.Event.GuildID == b.guildID && data.Event.ChannelID == b.channelID) {
		return false
	}

	userID := data.Event.SenderID()
	for _, admin := range b.superAdmins {
		if admin == userID {
			return true
		}
	}
	return false
}

func ErrAccessForbidden() *api.InteractionResponseData {
	return errorResponse(fmt.Errorf("access forbidden"))
}

func errorResponse(err error) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content:         option.NewNullableString("**Error:** " + err.Error()),
		Flags:           discord.EphemeralMessage,
		AllowedMentions: &api.AllowedMentions{ /* none */ },
	}
}

func (b *Bot) TxQueries(ctx context.Context) (q *sqlc.Queries, closer func(error) error, err error) {
	tx, closer, err := b.db.Tx(ctx)
	if err != nil {
		return nil, nil, err
	}
	return sqlc.New(tx), closer, nil
}

func (b *Bot) ConnQueries(ctx context.Context) (q *sqlc.Queries, closer func(), err error) {
	c, f, err := b.db.Conn(ctx)
	if err != nil {
		return nil, nil, err
	}
	return sqlc.New(c), f, nil
}

func (b *Bot) syncDatabaseState(ctx context.Context) (err error) {
	queries, closer, err := b.TxQueries(ctx)
	defer func() {
		err = closer(err)
	}()

	err = dao.RemovePlayerCountNotifications(ctx, queries)
	if err != nil {
		return err
	}

	trackings, err := dao.ListAllTrackings(ctx, queries)
	if err != nil {
		return err
	}

	//msgs := make([]*discord.Message, 0, len(trackings))
	notifications := make(map[model.MessageUserTarget]model.PlayerCountNotification)

	for _, t := range trackings {
		log.Printf("fetching message %s for notification tracking", t.MessageTarget)
		m, err := b.state.Message(t.ChannelID, t.MessageID)
		if err != nil {
			if ErrIsNotFound(err) || ErrIsAccessDenied(err) {
				// remove tracking of messages that were removed during downtime.
				err = dao.RemoveTrackingByMessageID(ctx, queries, t.GuildID, t.MessageID)
				if err != nil {
					return err
				}
				continue
			}
			return err
		}

		// iterate over all message reactions
		for _, reaction := range m.Reactions {
			emoji := reaction.Emoji.APIString()
			if _, ok := reactionPlayerCountNotificationMap[emoji]; !ok {
				// none of the ones that we want to look at
				continue
			}
			log.Printf("fetching users for emoji %s of message %s", emoji, t.MessageTarget)
			users, err := b.state.Reactions(m.ChannelID, t.MessageID, emoji, 0)
			if err != nil {
				if ErrIsNotFound(err) || ErrIsAccessDenied(err) {
					continue
				}
				return err
			}
			val := reactionPlayerCountNotificationMap[emoji]

			log.Printf("found %d users for emoji %s of message %s", len(users), emoji, t.MessageTarget)
			for _, user := range users {
				userTarget := model.MessageUserTarget{
					MessageTarget: t.MessageTarget,
					UserID:        user.ID,
				}
				if n, ok := notifications[userTarget]; ok {
					// only persist the smallest threshold
					if val < n.Threshold {
						// remove previous reaction that has a bigger value
						err = b.state.DeleteUserReaction(
							n.ChannelID,
							n.MessageID,
							n.UserID,
							reactionPlayerCountNotificationReverseMap[n.Threshold],
						)
						if err != nil {
							return err
						}

						// update to new lower value
						n.Threshold = val
						notifications[userTarget] = n
					} else {
						// remove previous reaction that has a bigger value
						err = b.state.DeleteUserReaction(
							n.ChannelID,
							n.MessageID,
							n.UserID,
							reactionPlayerCountNotificationReverseMap[val],
						)
						if err != nil {
							return err
						}
					}
				} else {
					notifications[userTarget] = model.PlayerCountNotification{
						MessageUserTarget: userTarget,
						Threshold:         val,
					}
				}
			}
		}
	}

	values := utils.Values(notifications)
	sort.Sort(model.ByPlayerCountNotificationIDs(values))

	err = dao.SetPlayerCountNotifications(ctx, queries, values)
	if err != nil {
		return err
	}

	return nil
}
