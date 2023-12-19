package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Config struct {
	DiscordToken       string `koanf:"discord.token" short:"t" description:"Discord App token."`
	DiscordChannelID   string `koanf:"discord.channel.id" short:"i" description:"Discord Channel ID."`
	DiscordSuperAdmins string `koanf:"super.admins" short:"a" description:"Comma separated list of Discord User IDs that are super admins."`
	SuperAdmins        []discord.UserID

	DatabaseDir string `koanf:"dir" short:"d" description:"Database directory"`

	GuildIDString string `koanf:"discord.guild.id" short:"g" description:"Discord Guild ID"`
	GuildID       discord.GuildID
}

func (c *Config) Validate() error {
	if c.DiscordToken == "" {
		return errors.New("discord token is required")
	}

	if c.DiscordChannelID == "" {
		return errors.New("discord channel id is required")
	}

	if c.GuildIDString != "" {
		snowflake, err := discord.ParseSnowflake(c.GuildIDString)
		if err != nil {
			return fmt.Errorf("invalid guild id: %s: %w", c.GuildID, err)
		}
		c.GuildID = discord.GuildID(snowflake)
	} else {
		c.GuildID = discord.NullGuildID
	}

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

	_, err := os.Stat(c.DatabaseDir)
	if err != nil {
		return err
	}

	return nil
}
