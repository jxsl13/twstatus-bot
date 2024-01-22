package model

import "github.com/diamondburned/arikawa/v3/discord"

var (
	ReactionPlayerCountNotificationMap = map[discord.APIEmoji]int{
		discord.NewAPIEmoji(0, "1️⃣"): 1,
		discord.NewAPIEmoji(0, "2️⃣"): 2,
		discord.NewAPIEmoji(0, "3️⃣"): 3,
		discord.NewAPIEmoji(0, "4️⃣"): 4,
		discord.NewAPIEmoji(0, "5️⃣"): 5,
		discord.NewAPIEmoji(0, "6️⃣"): 6,
		discord.NewAPIEmoji(0, "7️⃣"): 7,
		discord.NewAPIEmoji(0, "8️⃣"): 8,
		discord.NewAPIEmoji(0, "9️⃣"): 9,
		discord.NewAPIEmoji(0, "🔟"):   10,
	}
	ReactionPlayerCountNotificationReverseMap = map[int]discord.APIEmoji{
		1:  discord.NewAPIEmoji(0, "1️⃣"),
		2:  discord.NewAPIEmoji(0, "2️⃣"),
		3:  discord.NewAPIEmoji(0, "3️⃣"),
		4:  discord.NewAPIEmoji(0, "4️⃣"),
		5:  discord.NewAPIEmoji(0, "5️⃣"),
		6:  discord.NewAPIEmoji(0, "6️⃣"),
		7:  discord.NewAPIEmoji(0, "7️⃣"),
		8:  discord.NewAPIEmoji(0, "8️⃣"),
		9:  discord.NewAPIEmoji(0, "9️⃣"),
		10: discord.NewAPIEmoji(0, "🔟"),
	}
)
