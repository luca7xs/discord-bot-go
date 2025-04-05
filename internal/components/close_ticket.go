package components

import (
	"discord-bot-go/internal/db"
	"discord-bot-go/internal/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	// Botão "Fechar Ticket" - Exibe o modal de confirmação
	RegisterComponent("close_ticket_", func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		userID := i.Member.User.ID
		channelID := strings.TrimPrefix(i.MessageComponentData().CustomID, "close_ticket_")

		// Buscar o ticket para verificar permissões
		ticket, err := db.FindTicketByChannelID(channelID)
		if err != nil {
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Title:       "Erro",
				Description: "Erro ao verificar o ticket. Tente novamente!",
				Color:       0xFF0000, // Vermelho para erro
			})
			return
		}
		if ticket == nil {
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Title:       "Ticket Não Encontrado",
				Description: "Este ticket não foi encontrado. Ele pode já ter sido fechado!",
				Color:       0xFFFF00, // Amarelo para aviso
			})
			return
		}
		if ticket.UserID != userID && !utils.HasAdminPermission(s, i.Member, i.GuildID) {
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Title:       "Sem Permissão",
				Description: "Você não tem permissão para fechar este ticket!",
			})
			return
		}

		// Exibir modal de confirmação
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "confirm_close_" + channelID,
				Title:    "Confirmar Fechamento do Ticket",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "confirmation",
								Label:       "Digite 'sim' ou 'não' para confirmar.",
								Style:       discordgo.TextInputShort,
								Placeholder: "sim",
								Required:    true,
								MaxLength:   3,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "save_logs",
								Label:       "Salvar logs? (sim/não)",
								Style:       discordgo.TextInputShort,
								Placeholder: "sim",
								Required:    true,
								MaxLength:   3,
							},
						},
					},
				},
			},
		})
		if err != nil {
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Title:       "Erro",
				Description: "Erro ao exibir o formulário de confirmação. Tente novamente!",
				Color:       0xFF0000, // Vermelho para erro
			})
		}
	})
}
