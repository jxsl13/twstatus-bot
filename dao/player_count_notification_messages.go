package dao

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
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
