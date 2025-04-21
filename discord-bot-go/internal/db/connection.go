package db

import (
	"discord-bot-go/internal/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB é a conexão global com o banco de dados
var DB *gorm.DB

// InitDB inicializa a conexão com o banco de dados
func InitDB() {
	dsn := config.DatabaseURL
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Migrar os modelos
	err = DB.AutoMigrate(&Ticket{}, &TicketMessage{})
	if err != nil {
		log.Fatalf("Erro ao migrar o banco de dados: %v", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso!")
}
