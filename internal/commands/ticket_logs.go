package commands

import (
	"discord-bot-go/internal/db"
	"discord-bot-go/internal/utils"
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "ticket_logs",
			Description: "Obtenha os logs de um ticket fechado",
			DefaultMemberPermissions: func() *int64 {
				perm := int64(discordgo.PermissionAdministrator) // Apenas administradores podem usar
				return &perm
			}(),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "ticket",
					Description:  "Selecione um ticket fechado para ver os logs",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "de",
					Description: "Data inicial (formato: DD/MM/AAAA)",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "ate",
					Description: "Data final (formato: DD/MM/AAAA)",
					Required:    false,
				},
			},
		},
		Handler:             handleTicketLogsCommand,
		AutocompleteHandler: handleTicketLogsAutocomplete,
	})
}

func handleTicketLogsAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	tickets, err := db.FindClosedTicketsWithLogs()
	if err != nil {
		return
	}

	caser := cases.Title(language.BrazilianPortuguese)

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, t := range tickets {
		date := t.CreatedAt.Format("02/01")
		label := fmt.Sprintf("[#%d] %s - %s - %s", t.ID, t.UserName, date, caser.String(t.Type))

		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  label,
			Value: fmt.Sprint(t.ID),
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func handleTicketLogsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	opts := i.ApplicationCommandData().Options

	ticketID := opts[0].StringValue()
	ticket, err := db.FindTicketByID(ticketID)
	if err != nil || ticket == nil {
		utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
			Title:       "Ticket não encontrado",
			Description: "Verifique o ID do ticket selecionado.",
			Color:       0xFF0000,
		})
		return
	}

	messages, err := db.FindMessagesByTicketID(ticket.ID)
	if err != nil {
		utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
			Title:       "Erro ao buscar mensagens",
			Description: err.Error(),
			Color:       0xFF0000,
		})
		return
	}

	if len(messages) == 0 {
		utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
			Title:       "Logs não encontrados",
			Description: "Este ticket não possui logs de mensagens.",
			Color:       0xFFFF00,
		})
		return
	}

	// Parse de datas opcionais
	var startDate, endDate time.Time
	layout := "02/01/2006"

	for _, opt := range opts {
		switch opt.Name {
		case "de":
			startDate, _ = time.Parse(layout, opt.StringValue())
		case "ate":
			endDate, _ = time.Parse(layout, opt.StringValue())
		}
	}

	// Filtro de período
	var filteredMessages []*db.TicketMessage
	for i := range messages {
		msg := &messages[i]
		if (startDate.IsZero() || !msg.Timestamp.Before(startDate)) &&
			(endDate.IsZero() || !msg.Timestamp.After(endDate)) {
			filteredMessages = append(filteredMessages, msg)
		}
	}

	if len(filteredMessages) == 0 && (!startDate.IsZero() || !endDate.IsZero()) {
		utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
			Title:       "Nenhuma mensagem no período",
			Description: "Não foram encontradas mensagens dentro do período especificado.",
			Color:       0xFFFF00,
		})
		return
	}

	// Mensagem de feedback inicial (adiada)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Gerando o PDF com os logs do ticket...",
		},
	})

	pdfBytes, err := utils.GenerateTicketPDF(ticket, filteredMessages, s)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Erro ao gerar PDF",
					Description: err.Error(),
					Color:       0xFF0000,
				},
			},
		})
		return
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Files: []*discordgo.File{
			{
				Name:        fmt.Sprintf("ticket_%d_logs.pdf", ticket.ID),
				ContentType: "application/pdf",
				Reader:      strings.NewReader(string(pdfBytes)),
			},
		},
	})
}
