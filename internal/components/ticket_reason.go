package components

import (
	"discord-bot-go/internal/db"
	"discord-bot-go/internal/utils"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	// Registro do componente para o menu de seleção de tipo de ticket
	RegisterComponent("ticket_reason_", func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Pegar o tipo de ticket selecionado
		userID := i.Member.User.ID
		guildID := i.GuildID

		reason := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		ticketType := strings.Split(i.ModalSubmitData().CustomID, "_")[3]

		// Criar um canal privado para o ticket
		channelName := fmt.Sprintf("ticket-%s-%s", ticketType, userID[:5])
		channel, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
			Name: channelName,
			Type: discordgo.ChannelTypeGuildText,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:    userID, // Permissão para o usuário que criou o ticket
					Type:  discordgo.PermissionOverwriteTypeMember,
					Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
				},
				{
					ID:    s.State.User.ID, // Permissão para o bot
					Type:  discordgo.PermissionOverwriteTypeMember,
					Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
				},
				{
					ID:   guildID, // Bloquear @everyone
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel,
				},
			},
		})
		if err != nil {
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Description: "Erro ao criar o canal do ticket. Tente novamente!",
				Color:       0xFF0000, // Vermelho para erro
			})
			return
		}

		// Criar o ticket no banco de dados
		err = db.CreateTicket(userID, channel.ID, ticketType, reason)
		if err != nil {
			s.ChannelDelete(channel.ID)
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Description: "Erro ao registrar o ticket. Tente novamente!",
				Color:       0xFF0000, // Vermelho para erro
			})
			return
		}

		// Enviar uma mensagem inicial no canal do ticket
		welcomeMessage := fmt.Sprintf("Olá <@%s>! Seu ticket de **%s** foi criado. Motivo: **%s** Por favor, descreva sua solicitação.", userID, ticketType, reason)
		_, err = s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
			Content: welcomeMessage,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Fechar Ticket",
							Style:    discordgo.DangerButton,
							CustomID: "close_ticket_" + channel.ID,
						},
					},
				},
			},
		})
		if err != nil {
			s.ChannelDelete(channel.ID) // Deletar o canal se falhar ao enviar a mensagem
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Description: "Erro ao configurar o ticket. Tente novamente!",
				Color:       0xFF0000, // Vermelho para erro
			})
			return
		}

		// Responder ao usuário de forma efêmera
		utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
			Description: fmt.Sprintf("Ticket criado com sucesso! Veja o canal: <#%s>", channel.ID),
		})
	})
}
