package config

import (
	"errors"
	"fmt"
	"net/netip"
	"os"
	"strings"
)

type Config struct {
	DiscordToken       string `koanf:"token" short:"t" description:"Discord App token."`
	DiscordChannelID   string `koanf:"channel.id" short:"i" description:"Discord Channel ID."`
	DiscordSuperAdmins string `koanf:"super.admins" short:"a" description:"Comma separated list of Discord User IDs that are super admins."`
	admins             []string

	DatabaseDir string `koanf:"dir" short:"d" description:"Database directory"`

	TeeworldsServers string `koanf:"servers" short:"s" description:"Comma separated list of server addresses ip:port"`
	servers          []netip.AddrPort
}

func (c *Config) Validate() error {
	if c.DiscordToken == "" {
		return errors.New("discord token is required")
	}

	if c.DiscordChannelID == "" {
		return errors.New("discord channel id is required")
	}

	if c.DiscordSuperAdmins == "" {
		return errors.New("discord super admins is required")
	}
	c.admins = strings.Split(c.DiscordSuperAdmins, ",")
	if len(c.admins) == 0 {
		return errors.New("at least one discord super admin is required")
	}

	if c.DatabaseDir == "" {
		return errors.New("database directory is required")
	}

	_, err := os.Stat(c.DatabaseDir)
	if err != nil {
		return err
	}

	if c.TeeworldsServers == "" {
		return errors.New("teeworlds servers is required")
	}

	servers := []netip.AddrPort{}
	for _, addr := range strings.Split(c.TeeworldsServers, ",") {
		addrPort, err := netip.ParseAddrPort(addr)
		if err != nil {
			return fmt.Errorf("invalid teeworlds server address: %s: %w", addr, err)
		}
		servers = append(servers, addrPort)
	}

	c.servers = servers

	return nil
}

func (c *Config) Servers() []netip.AddrPort {
	return c.servers
}
