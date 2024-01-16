package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jxsl13/twstatus-bot/bot"
	"github.com/jxsl13/twstatus-bot/config"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/db"
	"github.com/spf13/cobra"
)

func main() {
	err := NewRootCmd().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	rootContext := rootContext{Ctx: ctx}

	// cmd represents the run command
	cmd := &cobra.Command{
		Use:   "twstatus-bot",
		Short: "twstatus-bot is a server status bot for Discord",
		RunE:  rootContext.RunE,
		Args:  cobra.ExactArgs(0),
		PostRunE: func(cmd *cobra.Command, args []string) error {

			cancel()
			return rootContext.DB.Close()
		},
	}

	// register flags but defer parsing and validation of the final values
	cmd.PreRunE = rootContext.PreRunE(cmd)

	// register flags but defer parsing and validation of the final values
	cmd.AddCommand(NewCompletionCmd(cmd.Name()))
	return cmd
}

type rootContext struct {
	Ctx    context.Context
	Config *config.Config
	DB     *db.DB
}

func (c *rootContext) PreRunE(cmd *cobra.Command) func(cmd *cobra.Command, args []string) error {

	c.Config = &config.Config{
		PostgresHostname: "postgres",
		PostgresPort:     5432,
		PostgresSSL:      false,
		PostgresDatabase: "twdb",
		PollInterval:     16 * time.Second,
	}
	runParser := config.RegisterFlags(c.Config, true, cmd)
	return func(cmd *cobra.Command, args []string) error {
		err := runParser()
		if err != nil {
			return err
		}

		err = dao.InitDatabase(c.Ctx, c.DB)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}

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
