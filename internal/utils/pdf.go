package utils

import (
	"fmt"
	"strings"

	"discord-bot-go/internal/db"

	"github.com/jung-kurt/gofpdf"
)

func GenerateTicketPDF(ticket *db.Ticket, messages []*db.TicketMessage) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Cabeçalho
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Ticket #%d - %s", ticket.ID, strings.ToUpper(ticket.Status)))
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Usuário: %s (%s)", ticket.UserName, ticket.UserID))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Tipo: %s", ticket.Type))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Motivo: %s", ticket.Reason))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Criado em: %s", ticket.CreatedAt.Format("02/01/2006 15:04")))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Fechado em: %s", ticket.ClosedAt.Format("02/01/2006 15:04")))
	pdf.Ln(12)

	// Mensagens
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Mensagens:")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	for _, m := range messages {
		author := m.UserName
		timestamp := m.Timestamp.Format("02/01 15:04")
		content := fmt.Sprintf("[%s] %s: %s", timestamp, author, m.Content)

		pdf.MultiCell(0, 6, content, "", "", false)
		pdf.Ln(2)
	}

	var buf strings.Builder
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}
