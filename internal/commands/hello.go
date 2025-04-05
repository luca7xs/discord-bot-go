package commands

import (
	"discord-bot-go/internal/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "hello",
			Description: "Diz olá para alguém!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "nome",
					Description: "Nome da pessoa a cumprimentar",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			nome := options[0].StringValue()
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Title:       "Olá!",
				Description: fmt.Sprintf("Olá, %s!", nome),
			})
		},
	})
}
