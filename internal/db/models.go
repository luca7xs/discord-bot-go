package db

import (
	"time"
)

// Ticket representa um ticket no banco de dados
type Ticket struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    string `gorm:"not null;index"`
	ChannelID string `gorm:"not null;index"`
	Type      string `gorm:"not null"`
	Reason    string `gorm:"not null"`
	Status    string `gorm:"not null;default:open;index"` // "open" ou "closed"
	CreatedAt time.Time
	ClosedAt  time.Time
	UpdatedAt time.Time
}

// TicketMessage representa uma mensagem de um ticket para os logs
type TicketMessage struct {
	ID        uint   `gorm:"primaryKey"`
	TicketID  uint   `gorm:"not null;index"` // Indexado para buscas por ticket
	UserID    string `gorm:"not null"`
	Content   string `gorm:"not null"`
	Timestamp time.Time
	CreatedAt time.Time
}
