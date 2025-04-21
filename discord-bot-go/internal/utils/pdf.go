package utils

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"discord-bot-go/internal/config"
	"discord-bot-go/internal/db"

	"github.com/bwmarrin/discordgo"
	"github.com/jung-kurt/gofpdf"
)

// formatDuration converts a time.Duration into a readable string like "2 dias, 4 horas, 15 minutos"
func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "0 segundos"
	}

	// Convert duration to days, hours, minutes, and seconds
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	// Build the parts of the duration string
	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", days, pluralize("dia", days)))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", hours, pluralize("hora", hours)))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", minutes, pluralize("minuto", minutes)))
	}
	if seconds > 0 || (days == 0 && hours == 0 && minutes == 0) {
		parts = append(parts, fmt.Sprintf("%d %s", seconds, pluralize("segundo", seconds)))
	}

	// Join parts with commas and handle the last part for proper formatting
	switch len(parts) {
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " e " + parts[1]
	default:
		return strings.Join(parts[:len(parts)-1], ", ") + " e " + parts[len(parts)-1]
	}
}

// pluralize returns the singular or plural form of a word based on the count
func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func GenerateTicketPDF(ticket *db.Ticket, messages []*db.TicketMessage, s *discordgo.Session) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Adiciona a fonte UTF8 - Mantenha isso!
	pdf.AddUTF8Font("Montserrat", "", "fonts/Montserrat-Regular.ttf")
	pdf.AddUTF8Font("Montserrat", "B", "fonts/Montserrat-Bold.ttf")

	pdf.AddPage()

	// Cabeçalho - Use a fonte Montserrat SEM o tradutor 'tr'
	pdf.SetFont("Montserrat", "B", 16)
	// Remova tr() daqui:
	pdf.Cell(40, 10, fmt.Sprintf("Ticket #%d - %s", ticket.ID, strings.ToUpper(ticket.Status)))
	pdf.Ln(12)

	pdf.SetFont("Montserrat", "", 12)
	// Usuário
	pdf.Cell(0, 10, fmt.Sprintf("Usuário: %s (%s)", ticket.UserName, ticket.UserID))
	pdf.Ln(8)
	// Tipo do ticket
	pdf.Cell(0, 10, fmt.Sprintf("Tipo: %s", ticket.Type))
	pdf.Ln(8)
	// Motivo de abertura
	pdf.Cell(0, 10, fmt.Sprintf("Motivo: %s", ticket.Reason))
	pdf.Ln(8)
	// Data de abertura
	pdf.Cell(0, 10, fmt.Sprintf("Criado em: %s", ticket.CreatedAt.Format("02/01/2006 15:04:05")))
	pdf.Ln(8)
	// Data de fechamento
	pdf.Cell(0, 10, fmt.Sprintf("Fechado em: %s", ticket.ClosedAt.Format("02/01/2006 15:04:05")))
	pdf.Ln(8)
	// Tempo de abertura
	pdf.Cell(0, 10, fmt.Sprintf("Tempo aberto: %s", formatDuration(ticket.ClosedAt.Sub(ticket.CreatedAt))))
	pdf.Ln(12)

	// Mensagens - Use a fonte Montserrat SEM o tradutor 'tr'
	pdf.SetFont("Montserrat", "B", 14)
	// Remova tr() daqui:
	pdf.Cell(0, 10, "Mensagens:")
	pdf.Ln(10)

	pdf.SetFont("Montserrat", "", 11)
	for _, m := range messages {
		author := m.UserName
		if author == "" {
			author = m.UserID
		}

		// Ignorar mensagens sem conteúdo
		if m.Content == "" {
			continue
		}

		// Determinar o tipo de usuário
		userType := "Tipo de usuário desconhecido"
		if m.UserID == ticket.UserID {
			userType = "Dono"
		} else if bot := s.State.User; bot != nil && bot.ID == m.UserID {
			userType = "Bot"
		} else {
			// Verificar permissões de administrador
			member, err := s.GuildMember(config.GuildID, m.UserID)
			if err == nil && member != nil {
				for _, roleID := range member.Roles {
					role, err := s.State.Role(config.GuildID, roleID)
					if err == nil && role.Permissions&discordgo.PermissionAdministrator != 0 {
						userType = "Administrador"
						break
					}
				}
			}
		}

		// Definir a cor apenas para o cabeçalho da mensagem
		switch userType {
		case "Dono":
			pdf.SetTextColor(0, 128, 0) // Verde
		case "Bot":
			pdf.SetTextColor(0, 0, 255) // Azul
		case "Administrador":
			pdf.SetTextColor(255, 0, 0) // Vermelho
		}

		timestamp := m.Timestamp.Format("02/01 15:04:05") // Adicionar segundos no timestamp
		header := fmt.Sprintf("[%s] %s (%s):", timestamp, author, userType)

		// Adicionar cabeçalho com cor
		pdf.MultiCell(0, 6, header, "", "", false)

		content := removeEmojis(m.Content)

		// Resetar a cor para preto e adicionar o conteúdo da mensagem
		pdf.SetTextColor(0, 0, 0) // Preto
		pdf.MultiCell(0, 6, content, "", "", false)
		pdf.Ln(2)
	}

	// Adicionar um indicador de fim das logs
	pdf.SetFont("Montserrat", "B", 12)
	pdf.SetTextColor(128, 128, 128) // Cinza
	pdf.Ln(10)
	pdf.Cell(0, 10, "--- Fim das Logs ---")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func removeEmojis(content string) string {
	// Regex para detectar emojis (simplificado)
	re := regexp.MustCompile(`[\x{1F000}-\x{1FFFF}]`)
	return re.ReplaceAllString(content, "[EMOJI]")
}
