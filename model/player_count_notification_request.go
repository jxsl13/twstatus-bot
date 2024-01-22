package model

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

type MessageUserTarget struct {
	MessageTarget
	UserID discord.UserID
}

type PlayerCountNotificationRequest struct {
	MessageUserTarget
	Threshold int
}

func (p *PlayerCountNotificationRequest) ToSetSQLC() sqlc.SetPlayerCountNotificationRequestParams {
	return sqlc.SetPlayerCountNotificationRequestParams{
		GuildID:   int64(p.GuildID),
		ChannelID: int64(p.ChannelID),
		MessageID: int64(p.MessageID),
		UserID:    int64(p.UserID),
		Threshold: int16(p.Threshold),
	}
}

func (p *PlayerCountNotificationRequest) ToRemoveSQLC() sqlc.RemovePlayerCountNotificationRequestParams {
	return sqlc.RemovePlayerCountNotificationRequestParams{
		GuildID:   int64(p.GuildID),
		ChannelID: int64(p.ChannelID),
		MessageID: int64(p.MessageID),
		UserID:    int64(p.UserID),
		Threshold: int16(p.Threshold),
	}
}

type PlayerCountNotificationRequests []PlayerCountNotificationRequest

type ByPlayerCountNotificationRequestIDs []PlayerCountNotificationRequest

func (a ByPlayerCountNotificationRequestIDs) Len() int      { return len(a) }
func (a ByPlayerCountNotificationRequestIDs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPlayerCountNotificationRequestIDs) Less(i, j int) bool {
	return a[i].MessageTarget.Less(a[j].MessageTarget)
}
