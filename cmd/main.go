package main

import (
	"discord-bot-go/internal/bot"
	"discord-bot-go/internal/db"
	"log"
)

func main() {
	// Inicializar o banco de dados
	db.InitDB()

	b, err := bot.NewBot()
	if err != nil {
		log.Fatalf("Erro ao criar bot: %v", err)
	}

	if err := b.Start(); err != nil {
		log.Fatalf("Erro ao iniciar bot: %v", err)
	}
}
