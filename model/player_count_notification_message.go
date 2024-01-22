package model

import (
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/sqlc"
	"github.com/jxsl13/twstatus-bot/utils"
)

func NewPlayerCountNotificationMessages(rows []sqlc.GetPlayerCountNotificationMessagesRow) []PlayerCountNotificationMessage {

	resultMap := make(map[ChannelTarget]PlayerCountNotificationMessage, len(rows)/10)
	for _, row := range rows {
		target := ChannelTarget{
			GuildID:   discord.GuildID(row.GuildID),
			ChannelID: discord.ChannelID(row.ChannelID),
		}

		usm := UserMessageThreshold{
			MessageThreshold: MessageThreshold{
				MessageID: discord.MessageID(row.ReqMessageID),
				Threshold: int(row.Threshold),
			},
			UserID: discord.UserID(row.UserID),
		}

		mt := MessageThreshold{
			MessageID: discord.MessageID(row.ReqMessageID),
			Threshold: int(row.Threshold),
		}

		if n, ok := resultMap[target]; ok {
			n.UserIDs = append(n.UserIDs, discord.UserID(row.UserID))
			n.RemoveMessageReactions = append(n.RemoveMessageReactions, mt)
			n.RemoveUserMessageReactions = append(n.RemoveUserMessageReactions, usm)
			resultMap[target] = n
		} else {
			resultMap[target] = PlayerCountNotificationMessage{
				ChannelTarget:              target,
				PrevMessageID:              discord.MessageID(row.PrevMessageID),
				UserIDs:                    []discord.UserID{discord.UserID(row.UserID)},
				RemoveMessageReactions:     []MessageThreshold{mt},
				RemoveUserMessageReactions: []UserMessageThreshold{usm},
			}
		}
	}

	result := make([]PlayerCountNotificationMessage, 0, len(resultMap))
	for _, v := range resultMap {
		v.UserIDs = utils.Unique(v.UserIDs)
		v.RemoveMessageReactions = utils.Unique(v.RemoveMessageReactions)
		result = append(result, v)
	}

	return result

}

type MessageUserID struct {
	MessageID discord.MessageID
	UserID    discord.UserID
}

type MessageThreshold struct {
	MessageID discord.MessageID
	// corresponding emoji must be removed
	// because we do not want to spam the channel
	Threshold int
}

func (m MessageThreshold) Reaction() discord.APIEmoji {
	return ReactionPlayerCountNotificationReverseMap[m.Threshold]
}

type UserMessageThreshold struct {
	MessageThreshold
	UserID discord.UserID
}

type PlayerCountNotificationMessage struct {
	// message is supposed to be sent into that channel
	ChannelTarget

	// needs to be removed from that channel if it exists
	// is 0 if no message was sent yet
	PrevMessageID discord.MessageID

	// mention these users for the current channel
	UserIDs []discord.UserID

	// for removing reactions from messages
	RemoveMessageReactions []MessageThreshold
	// for removing from database
	RemoveUserMessageReactions []UserMessageThreshold
}

func (p *PlayerCountNotificationMessage) MessageTarget(messageID discord.MessageID) MessageTarget {
	return MessageTarget{
		ChannelTarget: p.ChannelTarget,
		MessageID:     messageID,
	}
}

// format header
func (p *PlayerCountNotificationMessage) Format() string {

	const limit = 2000
	sb := strings.Builder{}
	sb.Grow(limit)

	for _, user := range p.UserIDs {
		mention := user.Mention()
		if sb.Len()+len(mention) > limit {
			break
		}
		sb.WriteString(mention)
		sb.WriteString(" ")
	}
	return sb.String()
}
