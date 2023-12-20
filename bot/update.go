package bot

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/servers"
)

func (b *Bot) updateServers(ctx context.Context) (src, dst int, err error) {
	servers, err := servers.GetAllServers()
	if err != nil {
		return 0, 0, err
	}

	serverList, err := model.NewServersFromDTO(servers)
	if err != nil {
		return 0, 0, err
	}

	tx, closer, err := b.Tx(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		err = closer(err)
	}()

	err = dao.SetServers(ctx, tx, serverList)
	if err != nil {
		return 0, 0, err
	}
	return len(servers), len(serverList), nil
}

func (b *Bot) updateServerList(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	if !b.IsSuperAdmin(data.Event.SenderID()) {
		return ErrAccessForbidden()
	}

	src, dst, err := b.updateServers(ctx)
	if err != nil {
		return errorResponse(err)
	}

	msg := fmt.Sprintf("Updated %d source to %d  target servers", src, dst)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}
