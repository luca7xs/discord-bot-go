package commands

import (
	"discord-bot-go/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Responde com Pong!",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Description: "Pong!",
			})
		},
	})
}
