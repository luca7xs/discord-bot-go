package components

import (
	"discord-bot-go/internal/config"

	"github.com/bwmarrin/discordgo"
)

func mapTicketTypesToSelectMenuOptions(ticketTypes []struct {
	Label       string
	Value       string
	Description string
}) []discordgo.SelectMenuOption {
	options := make([]discordgo.SelectMenuOption, len(ticketTypes))
	for i, t := range ticketTypes {
		options[i] = discordgo.SelectMenuOption{
			Label:       t.Label,
			Value:       t.Value,
			Description: t.Description,
		}
	}
	return options
}

func init() {
	RegisterComponent("create_ticket", func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Enviar uma mensagem efêmera com um SelectMenu para escolher o tipo de ticket
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Selecione o tipo de ticket:",
				Flags:   discordgo.MessageFlagsEphemeral,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID:    "ticket_type_" + i.Member.User.ID,
								Placeholder: "Escolha o tipo de ticket",
								Options:     mapTicketTypesToSelectMenuOptions(config.TicketTypes),
							},
						},
					},
				},
			},
		})
		if err != nil {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Erro ao enviar o menu de seleção. Tente novamente!",
			})
		}
	})
}
