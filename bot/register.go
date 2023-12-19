package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func (b *Bot) registerServer(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {

	return &api.InteractionResponseData{
		Content: option.NewNullableString("Server Registered"),
	}
}
