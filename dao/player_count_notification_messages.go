package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func (dao *DAO) AddPlayerCountNotificationMessage(ctx context.Context, channelID discord.ChannelID, messageID discord.MessageID) error {
	return dao.q.AddPlayerCountNotificationMessage(ctx, sqlc.AddPlayerCountNotificationMessageParams{
		ChannelID: int64(channelID),
		MessageID: int64(messageID),
	})
}

func (dao *DAO) RemovePlayerCountNotificationMessage(ctx context.Context, channelID discord.ChannelID, messageID discord.MessageID) error {
	return dao.q.RemovePlayerCountNotificationMessage(ctx, sqlc.RemovePlayerCountNotificationMessageParams{
		ChannelID: int64(channelID),
		MessageID: int64(messageID),
	})
}

func (dao *DAO) GetPlayerCountNotificationMessages(ctx context.Context, addresses []string) ([]model.PlayerCountNotificationMessage, error) {

	gpcnmr, err := dao.q.GetPlayerCountNotificationMessages(ctx, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to query player count notification messages: %w", err)
	}

	return model.NewPlayerCountNotificationMessages(gpcnmr), nil
}
