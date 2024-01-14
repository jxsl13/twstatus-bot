package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

// TODO: continue here
func ChangedMessageMentions(
	ctx context.Context,
	tx *sql.Tx,
	currentMentions model.MessageMentions,
) (
	messageMentions model.MessageMentions,
	err error,
) {

	// removed mentions
	// changed mentions
	// unchanged mentions

	return messageMentions, nil
}

func ListPrevMessageMentions(ctx context.Context, q *sqlc.Queries) (messageMentions model.MessageMentions, err error) {
	pmm, err := q.ListPreviousMessageMentions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query previous message mentions: %w", err)
	}
	messageMentions = make(model.MessageMentions, 128)
	var (
		mt model.MessageTarget
		u  discord.UserID
	)
	for _, pm := range pmm {
		mt = model.MessageTarget{
			ChannelTarget: model.ChannelTarget{
				GuildID:   discord.GuildID(pm.GuildID),
				ChannelID: discord.ChannelID(pm.ChannelID),
			},
			MessageID: discord.MessageID(pm.MessageID),
		}
		u = discord.UserID(pm.UserID)
		messageMentions[mt] = append(messageMentions[mt], u)
	}

	return messageMentions, nil
}

func RemoveMessageMentions(ctx context.Context, q *sqlc.Queries, mts []model.MessageTarget) (err error) {
	for _, mt := range mts {
		err = q.RemoveMessageMentions(ctx, sqlc.RemoveMessageMentionsParams{
			GuildID:   int64(mt.GuildID),
			ChannelID: int64(mt.ChannelID),
			MessageID: int64(mt.MessageID),
		})
		if err != nil {
			return fmt.Errorf("failed to remove message mentions: %w", err)
		}

	}
	return nil
}
