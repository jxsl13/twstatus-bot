package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func (dao *DAO) GetPlayerCountNotificationRequest(
	ctx context.Context,
	t model.MessageUserTarget,
) (
	notification model.PlayerCountNotificationRequest,
	err error,
) {

	ns, err := dao.q.GetPlayerCountNotificationRequest(ctx,
		sqlc.GetPlayerCountNotificationRequestParams{
			GuildID:   int64(t.GuildID),
			ChannelID: int64(t.ChannelID),
			MessageID: int64(t.MessageID),
			UserID:    int64(t.UserID),
		},
	)
	if err != nil {
		return model.PlayerCountNotificationRequest{}, err
	}
	if len(ns) == 0 {
		return model.PlayerCountNotificationRequest{}, fmt.Errorf("%w: player count notification", ErrNotFound)
	}
	n := ns[0]
	return model.PlayerCountNotificationRequest{
		MessageUserTarget: model.MessageUserTarget{
			UserID: discord.UserID(n.UserID),
			MessageTarget: model.MessageTarget{
				ChannelTarget: model.ChannelTarget{
					GuildID:   t.GuildID,
					ChannelID: t.ChannelID,
				},
				MessageID: t.MessageID,
			},
		},
		Threshold: int(n.Threshold),
	}, nil

}

func (dao *DAO) SetPlayerCountNotificationRequestList(ctx context.Context, notifications []model.PlayerCountNotificationRequest) (err error) {

	for _, n := range notifications {
		err = dao.q.SetPlayerCountNotificationRequest(ctx, sqlc.SetPlayerCountNotificationRequestParams{
			GuildID:   int64(n.GuildID),
			ChannelID: int64(n.ChannelID),
			MessageID: int64(n.MessageID),
			UserID:    int64(n.UserID),
			Threshold: int16(n.Threshold),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dao *DAO) SetPlayerCountNotificationRequest(ctx context.Context, n model.PlayerCountNotificationRequest) (err error) {
	return dao.q.SetPlayerCountNotificationRequest(ctx, n.ToSetSQLC())

}

func (dao *DAO) RemovePlayerCountNotificationRequests(ctx context.Context) (err error) {
	return dao.q.RemovePlayerCountNotificationRequests(ctx)
}

func (dao *DAO) RemovePlayerCountNotificationRequest(ctx context.Context, n model.PlayerCountNotificationRequest) (err error) {
	return dao.q.RemovePlayerCountNotificationRequest(ctx, n.ToRemoveSQLC())

}
