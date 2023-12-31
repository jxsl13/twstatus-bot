package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Config struct {
	DiscordToken       string `koanf:"discord.token" short:"t" description:"Discord App token."`
	DiscordSuperAdmins string `koanf:"super.admins" short:"a" description:"Comma separated list of Discord User IDs that are super admins."`
	SuperAdmins        []discord.UserID

	DatabaseDir string `koanf:"db.dir" short:"d" description:"Database directory"`
	WAL         bool   `koanf:"db.wal" short:"w" description:"Enable Write-Ahead-Log for SQLite"`

	GuildIDString string `koanf:"discord.guild.id" short:"g" description:"Discord Bot Owner Guild ID"`
	GuildID       discord.GuildID

	PollInterval        time.Duration `koanf:"poll.interval" short:"p" description:"Poll interval for DDNet's http master server"`
	LegacyMessageFormat bool          `koanf:"legacy.format" short:"l" description:"Use legacy message format. If disabled, rich text embeddings will be used."`
}

func (c *Config) Validate() error {
	if c.DiscordToken == "" {
		return errors.New("discord token is required")
	}

	snowflake, err := discord.ParseSnowflake(c.GuildIDString)
	if err != nil {
		return fmt.Errorf("invalid guild id: %s: %w", c.GuildID, err)
	}
	c.GuildID = discord.GuildID(snowflake)

	if c.DiscordSuperAdmins == "" {
		return errors.New("discord super admins is required")
	}
	admins := strings.Split(c.DiscordSuperAdmins, ",")
	if len(admins) == 0 {
		return errors.New("at least one discord super admin is required")
	} else {
		for _, admin := range admins {
			userID, err := discord.ParseSnowflake(admin)
			if err != nil {
				return fmt.Errorf("invalid discord super admin id: %s: %w", admin, err)
			}
			c.SuperAdmins = append(c.SuperAdmins, discord.UserID(userID))
		}
	}

	if c.DatabaseDir == "" {
		return errors.New("database directory is required")
	}

	_, err = os.Stat(c.DatabaseDir)
	if err != nil {
		return err
	}

	return nil
}
