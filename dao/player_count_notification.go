package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func (dao *DAO) ListAllPlayerCountNotifications(ctx context.Context) (notifications []model.PlayerCountNotification, err error) {
	pcn, err := dao.q.ListPlayerCountNotifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get player count notifications: %w", err)
	}
	result := make([]model.PlayerCountNotification, 0, len(pcn))
	for _, n := range pcn {
		result = append(result, model.PlayerCountNotification{
			MessageUserTarget: model.MessageUserTarget{
				UserID: discord.UserID(n.UserID),
				MessageTarget: model.MessageTarget{
					ChannelTarget: model.ChannelTarget{
						GuildID:   discord.GuildID(n.GuildID),
						ChannelID: discord.ChannelID(n.ChannelID),
					},
					MessageID: discord.MessageID(n.MessageID),
				},
			},
			Threshold: int(n.Threshold),
		})
	}
	return result, nil
}

func (dao *DAO) GetTargetListNotifications(
	ctx context.Context,
	servers map[model.MessageTarget]model.ChangedServerStatus) (
	_ map[model.MessageTarget]model.ChangedServerStatus,
	err error,
) {

	for t, server := range servers {

		messageNotifications, err := dao.q.GetMessageTargetNotifications(ctx,
			sqlc.GetMessageTargetNotificationsParams{
				GuildID:   int64(t.GuildID),
				ChannelID: int64(t.ChannelID),
				MessageID: int64(t.MessageID),
			})
		if err != nil {
			return nil, fmt.Errorf("failed to get message target notifications: %w", err)
		}

		for _, n := range messageNotifications {
			notification := model.PlayerCountNotification{
				Threshold: int(n.Threshold),
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
			}

			if notification.Notify(&server) {
				server.UserNotifications = append(server.UserNotifications, discord.UserID(n.UserID))
				servers[t] = server
			}
		}

	}
	return servers, nil
}

func (dao *DAO) GetPlayerCountNotification(
	ctx context.Context,
	t model.MessageUserTarget,
) (
	notification model.PlayerCountNotification,
	err error,
) {

	ns, err := dao.q.GetPlayerCountNotification(ctx,
		sqlc.GetPlayerCountNotificationParams{
			GuildID:   int64(t.GuildID),
			ChannelID: int64(t.ChannelID),
			MessageID: int64(t.MessageID),
			UserID:    int64(t.UserID),
		},
	)
	if err != nil {
		return model.PlayerCountNotification{}, err
	}
	if len(ns) == 0 {
		return model.PlayerCountNotification{}, fmt.Errorf("%w: player count notification", ErrNotFound)
	}
	n := ns[0]
	return model.PlayerCountNotification{
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

func (dao *DAO) SetPlayerCountNotificationList(ctx context.Context, notifications []model.PlayerCountNotification) (err error) {

	for _, n := range notifications {
		err = dao.q.SetPlayerCountNotification(ctx, sqlc.SetPlayerCountNotificationParams{
			GuildID:   int64(n.GuildID),
			ChannelID: int64(n.ChannelID),
			MessageID: int64(n.MessageID),
			UserID:    int64(n.UserID),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dao *DAO) SetPlayerCountNotification(ctx context.Context, n model.PlayerCountNotification) (err error) {
	return dao.q.SetPlayerCountNotification(ctx, n.ToSetSQLC())

}

func (dao *DAO) RemovePlayerCountNotifications(ctx context.Context) (err error) {
	return dao.q.RemovePlayerCountNotifications(ctx)
}

func (dao *DAO) RemovePlayerCountNotification(ctx context.Context, n model.PlayerCountNotification) (err error) {
	return dao.q.RemovePlayerCountNotification(ctx, n.ToRemoveSQLC())

}
