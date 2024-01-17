package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"time"

	"github.com/jxsl13/twstatus-bot/bot"
	"github.com/jxsl13/twstatus-bot/config"
	"github.com/jxsl13/twstatus-bot/db"
	"github.com/jxsl13/twstatus-bot/migrations"
	"github.com/spf13/cobra"
)

func main() {

	err := NewRootCmd(migrations.FS).Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewRootCmd(migrationsFs fs.FS) *cobra.Command {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	rootContext := rootContext{
		Ctx:          ctx,
		MigrationsFS: migrationsFs,
	}

	// cmd represents the run command
	cmd := &cobra.Command{
		Use:   "twstatus-bot",
		Short: "twstatus-bot is a server status bot for Discord",
		RunE:  rootContext.RunE,
		Args:  cobra.ExactArgs(0),
		PostRunE: func(cmd *cobra.Command, args []string) error {
			rootContext.DB.Close()

			cancel()
			return nil
		},
	}

	// register flags but defer parsing and validation of the final values
	cmd.PreRunE = rootContext.PreRunE(cmd)

	// register flags but defer parsing and validation of the final values
	cmd.AddCommand(NewCompletionCmd(cmd.Name()))
	return cmd
}

type rootContext struct {
	Ctx          context.Context
	MigrationsFS fs.FS

	// set in PreRunE
	Config *config.Config
	DB     *db.DB
}

func (c *rootContext) PreRunE(cmd *cobra.Command) func(cmd *cobra.Command, args []string) error {

	c.Config = &config.Config{
		PostgresHostname: "postgres",
		PostgresPort:     5432,
		PostgresSSLMode:  db.SSLModeDisable,
		PostgresDatabase: "twdb",
		PollInterval:     16 * time.Second,
	}
	runParser := config.RegisterFlags(c.Config, true, cmd)
	return func(cmd *cobra.Command, args []string) error {
		err := runParser()
		if err != nil {
			return err
		}

		db, err := db.New(
			c.Ctx,
			c.Config.PostgresHostname,
			c.Config.PostgresPort,
			c.Config.PostgresDatabase,
			c.Config.PostgresUser,
			c.Config.PostgresPassword,
			db.WithMigrationsFs(c.MigrationsFS),
			db.WithSSL(c.Config.PostgresSSLMode),
		)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		c.DB = db

		return nil
	}
}

func (c *rootContext) RunE(cmd *cobra.Command, args []string) error {
	b, err := bot.New(
		c.Ctx,
		c.Config.DiscordToken,
		c.DB,
		c.Config.SuperAdmins,
		c.Config.GuildID,
		c.Config.ChannelID,
		c.Config.PollInterval,
		c.Config.LegacyMessageFormat,
	)
	if err != nil {
		return err
	}
	defer b.Close()

	err = b.Connect(c.Ctx)
	if err != nil {
		return err
	}
	return nil
}
