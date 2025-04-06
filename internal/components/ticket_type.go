package components

import (
	"discord-bot-go/internal/db"
	"discord-bot-go/internal/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	// Registro do componente para o menu de seleção de tipo de ticket
	RegisterComponent(Component{
		CustomID: "ticket_type_",
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Pegar o tipo de ticket selecionado
			userID := i.Member.User.ID
			selectedType := i.MessageComponentData().Values[0]

			// Verificar se o usuário já tem um ticket aberto
			existingTicket, err := db.FindOpenTicketByUserIDAndType(userID, selectedType)
			if err != nil {
				utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
					Description: "Ocorreu um erro ao verificar se você possui tickets abertos. Tente novamente mais tarde.",
					Color:       0xFF0000, // Vermelho para erro
				})
				return
			}

			if existingTicket != nil {
				utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
					Description: fmt.Sprintf("Você já tem um ticket de **%s** aberto em <#%s>. Feche-o antes de criar outro!", selectedType, existingTicket.ChannelID),
					Color:       0xFFFF00, // Amarelo para aviso
				})
				return
			}

			// Enviar um modal para coletar o motivo do ticket
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "ticket_reason_" + userID + "_" + selectedType,
					Title:    "Motivo do Ticket",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "reason",
									Label:       "Por que você está abrindo este ticket?",
									Style:       discordgo.TextInputParagraph,
									Placeholder: "Descreva seu problema ou sugestão aqui...",
									Required:    true,
									MaxLength:   1000,
								},
							},
						},
					},
				},
			})
			if err != nil {
				utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
					Description: "Erro ao enviar o formulário. Tente novamente!",
					Color:       0xFF0000, // Vermelho para erro
				})
			}
		},
	})
}
