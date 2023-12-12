package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	"github.com/jxsl13/twstatus-bot/bot"
	"github.com/jxsl13/twstatus-bot/config"
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
	Ctx    context.Context
	Config *config.Config
	DB     *sql.DB
}

func (c *rootContext) PreRunE(cmd *cobra.Command) func(cmd *cobra.Command, args []string) error {
	c.Config = &config.Config{
		DatabaseDir: "./",
	}
	runParser := config.RegisterFlags(c.Config, true, cmd)
	return func(cmd *cobra.Command, args []string) error {
		err := runParser()
		if err != nil {
			return err
		}

		return nil
	}
}

func (c *rootContext) RunE(cmd *cobra.Command, args []string) error {
	b, err := bot.New(c.Config.DiscordToken)
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
