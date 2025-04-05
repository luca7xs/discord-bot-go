package commands

import (
	"discord-bot-go/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "info",
			Description: "Mostra informações básicas",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Title:       "Informações do Bot",
				Description: "Este é um bot feito em Go!",
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Linguagem", Value: "Golang", Inline: true},
					{Name: "Criador", Value: "Você!", Inline: true},
				},
			})
		},
	})
}
