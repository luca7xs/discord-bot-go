package components

import (
	"discord-bot-go/internal/db"
	"discord-bot-go/internal/utils"
	"fmt"
	"strings"
	"time"

	"slices"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterComponent(Component{
		CustomID: "confirm_close",
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channelID := strings.TrimPrefix(i.ModalSubmitData().CustomID, "confirm_close_")

			// Validar entradas do modal
			confirmation, err := validateInput(i.ModalSubmitData().Components[0], "sim", "não")
			if err != nil {
				utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
					Title:       "Erro de Entrada",
					Description: err.Error(),
					Color:       0xFF0000, // Vermelho para erro
				})
				return
			}

			saveLogs, err := validateInput(i.ModalSubmitData().Components[1], "sim", "não")
			if err != nil {
				utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
					Title:       "Erro de Entrada",
					Description: err.Error(),
					Color:       0xFF0000, // Vermelho para erro
				})
				return
			}

			// Cancelar se o usuário não confirmou
			if confirmation == "não" {
				utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
					Title:       "Cancelado",
					Description: "O fechamento do ticket foi cancelado.",
					Color:       0xFFFF00, // Amarelo para cancelamento
				})
				return
			}

			// Buscar o ticket para garantir que ainda existe
			ticket, err := db.FindTicketByChannelID(channelID)
			if err != nil || ticket == nil {
				utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
					Title:       "Erro",
					Description: "Ticket não encontrado durante o fechamento. Pode já ter sido fechado!",
					Color:       0xFF0000, // Vermelho para erro
				})
				return
			}

			// Responder ao modal para fechá-lo com sucesso
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Title:       "Confirmação Recebida",
				Description: "Aguarde o processamento...",
				Color:       0x00FF00, // Verde para sucesso
			})

			// Criar mensagem com botão de cancelamento
			saveLogsMsg := "com logs salvos"
			if saveLogs == "não" {
				saveLogsMsg = "sem salvar logs"
			}
			msg, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Ticket Fechado",
						Description: fmt.Sprintf("O ticket de **%s** será fechado %s! Esse canal será excluído em 20 segundos.", ticket.Type, saveLogsMsg),
						Color:       0x00FF00, // Verde para indicar que o processo está em andamento
					},
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Cancelar",
								Style:    discordgo.DangerButton,
								CustomID: "cancel_close_" + channelID,
							},
						},
					},
				},
			})
			if err != nil {
				s.ChannelMessageSend(channelID, "Erro ao enviar mensagem de fechamento. O canal será mantido ativo. Por favor, tente novamente ou contate um administrador.")
				return
			}

			// Gerenciar cancelamento e fechamento
			handleTicketClosure(s, channelID, msg, saveLogs == "sim", ticket)
		},
	})
}

// Função para validar entradas do modal
func validateInput(component discordgo.MessageComponent, validOptions ...string) (string, error) {
	input := strings.ToLower(component.(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
	if slices.Contains(validOptions, input) {
		return input, nil
	}
	return "", fmt.Errorf("por favor, digite apenas '%s'", strings.Join(validOptions, "' ou '"))
}

// Função para gerenciar o fechamento do ticket
func handleTicketClosure(s *discordgo.Session, channelID string, msg *discordgo.Message, saveLogs bool, ticket *db.Ticket) {
	cancelChan := make(chan bool, 1)

	// Registrar handler temporário para o botão de cancelamento
	cancelHandler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.MessageComponentData().CustomID == "cancel_close_"+channelID {
			cancelChan <- true
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Title:       "Fechamento Cancelado",
				Description: "O canal será mantido aberto.",
				Color:       0xFFFF00, // Amarelo para cancelamento
			})
			err := s.ChannelMessageDelete(msg.ChannelID, msg.ID)
			if err != nil {
				s.ChannelMessageSend(channelID, "Erro ao deletar mensagem de fechamento. O processo continuará mesmo assim.")
			}
		}
	}
	s.AddHandlerOnce(cancelHandler)

	// Aguardar 20 segundos ou cancelamento
	select {
	case <-cancelChan:
		return // Sai da função se cancelado, sem fechar o ticket no banco
	case <-time.After(20 * time.Second):
		// Fechar o ticket no banco de dados apenas se não for cancelado
		var messages []*discordgo.Message
		if saveLogs {
			var err error
			messages, err = utils.FetchAllMessages(s, channelID)
			if err != nil {
				s.ChannelMessageSend(channelID, "Erro ao buscar mensagens do ticket. O canal será mantido ativo. Por favor, tente novamente ou contate um administrador.")
				return
			}
		}

		err := db.CloseTicket(channelID, messages)
		if err != nil {
			s.ChannelMessageSend(channelID, "Erro ao registrar o fechamento do ticket no banco de dados. O canal será mantido ativo. Por favor, tente novamente ou contate um administrador.")
			return
		}

		_, err = s.ChannelDelete(channelID)
		if err != nil {
			s.ChannelMessageSend(channelID, "Ticket fechado no banco, mas erro ao deletar o canal. Por favor, delete manualmente!")
		}
	}
}
