package db

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// Constantes para status de tickets
const (
	TicketStatusOpen   = "open"
	TicketStatusClosed = "closed"
)

// CreateTicket cria um novo ticket no banco de dados
func CreateTicket(userID, userName, channelID, ticketType, reason string) error {
	if userID == "" || userName == "" || channelID == "" || ticketType == "" {
		return fmt.Errorf("userID, userName, channelID e ticketType não podem ser vazios")
	}

	ticket := Ticket{
		UserID:    userID,
		UserName:  userName,
		ChannelID: channelID,
		Type:      ticketType,
		Reason:    reason,
		Status:    TicketStatusOpen,
		CreatedAt: time.Now(),
	}
	if err := DB.Create(&ticket).Error; err != nil {
		return fmt.Errorf("falha ao criar ticket para userID %s: %w", userID, err)
	}
	return nil
}

// findTicket é uma função auxiliar genérica para buscar tickets
func findTicket(whereClause string, args ...interface{}) (*Ticket, error) {
	var ticket Ticket
	err := DB.Where(whereClause, args...).First(&ticket).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar ticket: %w", err)
	}
	return &ticket, nil
}

// FindOpenTicketByUserID busca um ticket aberto para um usuário
func FindOpenTicketByUserID(userID string) (*Ticket, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID não pode ser vazio")
	}
	return findTicket("user_id = ? AND status = ?", userID, TicketStatusOpen)
}

// FindOpenTicketByUserIDAndType busca um ticket aberto para um usuário e tipo específico
func FindOpenTicketByUserIDAndType(userID, ticketType string) (*Ticket, error) {
	if userID == "" || ticketType == "" {
		return nil, fmt.Errorf("userID e ticketType não podem ser vazios")
	}
	return findTicket("user_id = ? AND type = ? AND status = ?", userID, ticketType, TicketStatusOpen)
}

// FindTicketByChannelID busca um ticket pelo ChannelID
func FindTicketByChannelID(channelID string) (*Ticket, error) {
	if channelID == "" {
		return nil, fmt.Errorf("channelID não pode ser vazio")
	}
	return findTicket("channel_id = ?", channelID)
}

func FindClosedTicketsWithLogs() ([]Ticket, error) {
	var tickets []Ticket
	err := DB.
		Model(&Ticket{}).
		Joins("JOIN ticket_messages ON ticket_messages.ticket_id = tickets.id").
		Where("tickets.status = ?", TicketStatusClosed).
		Group("tickets.id").
		Find(&tickets).Error

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tickets fechados com logs: %w", err)
	}
	return tickets, nil
}

func FindTicketByID(id string) (*Ticket, error) {
	var ticket Ticket
	err := DB.First(&ticket, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func FindMessagesByTicketID(ticketID uint) ([]TicketMessage, error) {
	var messages []TicketMessage
	err := DB.Where("ticket_id = ?", ticketID).Order("timestamp asc").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// CloseTicket marca um ticket como fechado e salva as mensagens no banco de dados, se fornecidas
func CloseTicket(channelID string, messages []*discordgo.Message) error {
	if channelID == "" {
		return fmt.Errorf("channelID não pode ser vazio")
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		// Atualizar o ticket
		result := tx.Model(&Ticket{}).Where("channel_id = ? AND status = ?", channelID, TicketStatusOpen).Updates(map[string]interface{}{
			"status":    TicketStatusClosed,
			"closed_at": time.Now(),
		})
		if result.Error != nil {
			return fmt.Errorf("falha ao atualizar ticket com channelID %s: %w", channelID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("ticket com channelID %s não encontrado ou já fechado", channelID)
		}

		// Buscar o ticket para obter o ID
		var ticket Ticket
		if err := tx.Where("channel_id = ?", channelID).First(&ticket).Error; err != nil {
			return fmt.Errorf("falha ao buscar ticket com channelID %s: %w", channelID, err)
		}

		// Salvar as mensagens no banco de dados, se fornecidas
		for _, msg := range messages {
			ticketMsg := TicketMessage{
				TicketID:  ticket.ID,
				UserID:    msg.Author.ID,
				UserName:  msg.Author.Username,
				Content:   msg.Content,
				Timestamp: msg.Timestamp,
				CreatedAt: time.Now(),
			}
			if err := tx.Create(&ticketMsg).Error; err != nil {
				return fmt.Errorf("falha ao salvar mensagem para ticket %d: %w", ticket.ID, err)
			}
		}
		return nil
	})
}
