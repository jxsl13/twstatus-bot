package bot

import (
	"context"
	"fmt"
	"net/netip"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/model"
)

type AddTrackingParams struct {
	Address string `discord:"address"`
}

func (b *Bot) addTracking(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	var params AddTrackingParams
	err := data.Options.Unmarshal(&params)
	if err != nil {
		return errorResponse(err)
	}

	channelID := optionalChannelID(data)

	addresses := strings.Split(params.Address, ",")
	for _, address := range addresses {
		_, err = netip.ParseAddrPort(address)
		if err != nil {
			return errorResponse(fmt.Errorf("invalid address: %w", err))
		}
	}

	msgs := make([]*discord.Message, 0, len(addresses))
	defer func() {
		if err != nil {
			for idx, msg := range msgs {
				auditReason := fmt.Sprintf("failed to add tracking for %s", addresses[idx])
				_ = b.state.DeleteMessage(
					channelID,
					msg.ID,
					api.AuditLogReason(auditReason),
				)
			}
		}
	}()
	for _, address := range addresses {
		msg, err := b.state.SendMessage(channelID, fmt.Sprintf("initial message for %s tracking", address))
		if err != nil {
			return errorResponse(err)
		}
		msgs = append(msgs, msg)
	}

	q, closer, err := b.TxQueries(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	for idx, msg := range msgs {
		err = dao.AddTracking(ctx, q, model.Tracking{
			MessageTarget: model.MessageTarget{
				ChannelTarget: model.ChannelTarget{
					GuildID:   data.Event.GuildID,
					ChannelID: channelID,
				},
				MessageID: msg.ID,
			},
			Address: addresses[idx],
		})
		if err != nil {
			return errorResponse(err)
		}
	}

	plural := ""
	if len(addresses) != 1 {
		plural = "es"
	}

	msg := fmt.Sprintf("Added tracking for %d address%s", len(addresses), plural)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}
