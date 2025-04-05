package bot

import (
	"discord-bot-go/internal/commands"
	"discord-bot-go/internal/components"
	"discord-bot-go/internal/config"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
}

func NewBot() (*Bot, error) {
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar sessão: %v", err)
	}
	return &Bot{Session: dg}, nil
}

func (b *Bot) Start() error {
	// Handler dinâmico para interações (comandos e componentes)
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			// Lidar com comandos
			cmdName := i.ApplicationCommandData().Name
			if cmd, exists := commands.CommandRegistry[cmdName]; exists {
				cmd.Handler(s, i)
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			// Lidar com autocomplete
			cmdName := i.ApplicationCommandData().Name
			if cmd, exists := commands.CommandRegistry[cmdName]; exists {
				cmd.Handler(s, i)
			}
		case discordgo.InteractionMessageComponent:
			// Lidar com interações de componentes (botões, modals, etc.)
			customID := i.MessageComponentData().CustomID
			for prefix, component := range components.ComponentRegistry {
				if strings.HasPrefix(customID, prefix) {
					component.Handler(s, i)
					return
				}
			}
		case discordgo.InteractionModalSubmit:
			// Lidar com envios de modals
			customID := i.ModalSubmitData().CustomID
			for prefix, component := range components.ComponentRegistry {
				if strings.HasPrefix(customID, prefix) {
					component.Handler(s, i)
					return
				}
			}
		}
	})

	// Abre a conexão com o Discord
	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão: %v", err)
	}

	// Registra todos os comandos automaticamente
	err = b.registerCommands()
	if err != nil {
		return err
	}

	fmt.Println("Bot está rodando! Pressione CTRL+C para sair.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return b.Session.Close()
}

func (b *Bot) registerCommands() error {
	for _, cmd := range commands.CommandRegistry {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", cmd.Definition)
		if err != nil {
			return fmt.Errorf("erro ao registrar comando %s: %v", cmd.Definition.Name, err)
		}
		fmt.Printf("Comando /%s registrado com sucesso!\n", cmd.Definition.Name)
	}
	return nil
}
