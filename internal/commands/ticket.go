package commands

import (
	"discord-bot-go/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "ticket",
			Description: "Configura um sistema de tickets em um canal",
			DefaultMemberPermissions: func() *int64 {
				perm := int64(discordgo.PermissionAdministrator) // Apenas administradores podem usar
				return &perm
			}(),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "canal",
					Description: "Canal onde o sistema de tickets será configurado",
					Required:    true,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
					},
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Pegar o canal selecionado
			channelID := i.ApplicationCommandData().Options[0].ChannelValue(s).ID

			// Enviar a mensagem com o botão no canal especificado
			_, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Embed: &discordgo.MessageEmbed{
					Title:       "Sistema de Tickets",
					Description: "Clique no botão abaixo para criar um ticket de suporte!",
					Color:       0x2b2d31,
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Criar Ticket",
								Style:    discordgo.PrimaryButton,
								CustomID: "create_ticket",
							},
						},
					},
				},
			})
			if err != nil {
				utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
					Description: "Erro ao configurar o sistema de tickets. Tente novamente!",
				})
				return
			}

			// Responder ao administrador
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Description: "Sistema de tickets configurado com sucesso no canal <#" + channelID + ">!",
			})
		},
	})
}
