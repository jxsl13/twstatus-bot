package model

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

type MessageUserTarget struct {
	MessageTarget
	UserID discord.UserID
}

type PlayerCountNotification struct {
	MessageUserTarget
	Threshold int
}

func (p *PlayerCountNotification) ToSetSQLC() sqlc.SetPlayerCountNotificationParams {
	return sqlc.SetPlayerCountNotificationParams{
		GuildID:   int64(p.GuildID),
		ChannelID: int64(p.ChannelID),
		MessageID: int64(p.MessageID),
		UserID:    int64(p.UserID),
		Threshold: int16(p.Threshold),
	}
}

func (p *PlayerCountNotification) ToRemoveSQLC() sqlc.RemovePlayerCountNotificationParams {
	return sqlc.RemovePlayerCountNotificationParams{
		GuildID:   int64(p.GuildID),
		ChannelID: int64(p.ChannelID),
		MessageID: int64(p.MessageID),
		UserID:    int64(p.UserID),
	}
}

func (p PlayerCountNotification) Notify(change *ChangedServerStatus) bool {
	return len(change.Curr.Clients) >= p.Threshold
}

type PlayerCountNotifications []PlayerCountNotification

type ByPlayerCountNotificationIDs []PlayerCountNotification

func (a ByPlayerCountNotificationIDs) Len() int      { return len(a) }
func (a ByPlayerCountNotificationIDs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPlayerCountNotificationIDs) Less(i, j int) bool {
	return a[i].MessageTarget.Less(a[j].MessageTarget)
}
