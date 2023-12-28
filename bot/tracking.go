package bot

import (
	"context"
	"fmt"
	"net/netip"

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

	_, err = netip.ParseAddrPort(params.Address)
	if err != nil {
		return errorResponse(fmt.Errorf("invalid address: %w", err))
	}

	msg, err := b.state.SendMessage(channelID, fmt.Sprintf("initial message for %s tracking", params.Address))
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		if err != nil {
			auditReason := fmt.Sprintf("failed to add tracking for %s", params.Address)
			_ = b.state.DeleteMessage(
				channelID,
				msg.ID,
				api.AuditLogReason(auditReason),
			)
		}
	}()

	tx, closer, err := b.Tx(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	err = dao.AddTracking(ctx, tx, model.Tracking{
		GuildID:   data.Event.GuildID,
		ChannelID: channelID,
		Address:   params.Address,
		MessageID: msg.ID,
	})
	if err != nil {
		return errorResponse(err)
	}
	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("Added tracking for %s", params.Address)),
		Flags:   discord.EphemeralMessage,
	}
}
