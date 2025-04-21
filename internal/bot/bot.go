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
	// Abre a conexão com o Discord
	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão: %v", err)
	}

	fmt.Printf("Conectado ao Discord como %s#%s\n", b.Session.State.User.Username, b.Session.State.User.Discriminator)

	// Limpa todos os comandos registrados (globais e do servidor)
	// err = b.clearCommands()
	// if err != nil {
	// 	return fmt.Errorf("erro ao limpar comandos: %v", err)
	// }

	// Registra todos os comandos no servidor especificado
	err = b.registerCommands()
	if err != nil {
		return err
	}

	// Registra os handlers de interações
	err = b.registerInteractionsHandlers()
	if err != nil {
		return err
	}

	fmt.Println("Bot está rodando! Pressione CTRL+C para sair.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return b.Session.Close()
}

// clearCommands remove todos os comandos registrados (globais e do servidor)
// func (b *Bot) clearCommands() error {
// 	// Remove comandos globais
// 	globalCommands, err := b.Session.ApplicationCommands(b.Session.State.User.ID, "")
// 	if err != nil {
// 		return fmt.Errorf("erro ao buscar comandos globais: %v", err)
// 	}
// 	for _, cmd := range globalCommands {
// 		err = b.Session.ApplicationCommandDelete(b.Session.State.User.ID, "", cmd.ID)
// 		if err != nil {
// 			return fmt.Errorf("erro ao deletar comando global %s: %v", cmd.Name, err)
// 		}
// 		fmt.Printf("Comando global /%s deletado com sucesso!\n", cmd.Name)
// 	}

// 	// Remove comandos do servidor especificado
// 	guildCommands, err := b.Session.ApplicationCommands(b.Session.State.User.ID, config.GuildID)
// 	if err != nil {
// 		return fmt.Errorf("erro ao buscar comandos do servidor %s: %v", config.GuildID, err)
// 	}
// 	for _, cmd := range guildCommands {
// 		err = b.Session.ApplicationCommandDelete(b.Session.State.User.ID, config.GuildID, cmd.ID)
// 		if err != nil {
// 			return fmt.Errorf("erro ao deletar comando do servidor %s: %v", cmd.Name, err)
// 		}
// 		fmt.Printf("Comando do servidor /%s deletado com sucesso!\n", cmd.Name)
// 	}

// 	fmt.Println("Todos os comandos foram limpos com sucesso!")
// 	return nil
// }

// registerCommands registra todos os comandos no servidor especificado
func (b *Bot) registerCommands() error {
	for _, cmd := range commands.CommandRegistry {
		// Registra o comando no servidor especificado (config.GuildID)
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, config.GuildID, cmd.Definition)
		if err != nil {
			return fmt.Errorf("erro ao registrar comando %s: %v", cmd.Definition.Name, err)
		}
		fmt.Printf("Comando /%s registrado com sucesso no servidor %s!\n", cmd.Definition.Name, config.GuildID)
	}
	return nil
}

// registerInteractionsHandlers registra os handlers de interações para comandos, menus, botões, etc.
func (b *Bot) registerInteractionsHandlers() error {
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Verifica se a interação ocorreu no servidor correto
		if i.GuildID != config.GuildID {
			// Ignora interações de outros servidores
			fmt.Printf("Interação ignorada: recebida do servidor %s, esperado %s\n", i.GuildID, config.GuildID)
			return
		}

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			// Lidar com comandos
			cmdName := i.ApplicationCommandData().Name
			if cmd, exists := commands.CommandRegistry[cmdName]; exists {
				cmd.Handler(s, i)
				return
			}
		case discordgo.InteractionMessageComponent:
			// Lidar com interações de componentes (botões, modals, menus, etc.)
			customID := i.MessageComponentData().CustomID
			for prefix, component := range components.ComponentRegistry {
				if strings.HasPrefix(customID, prefix) {
					component.Handler(s, i)
					return
				}
			}
		case discordgo.InteractionModalSubmit:
			// Lidar com interações de modais
			customID := i.ModalSubmitData().CustomID
			for prefix, component := range components.ComponentRegistry {
				if strings.HasPrefix(customID, prefix) {
					component.Handler(s, i)
					return
				}
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			// Lidar com autocomplete
			cmdName := i.ApplicationCommandData().Name
			if cmd, exists := commands.CommandRegistry[cmdName]; exists && cmd.AutocompleteHandler != nil {
				cmd.AutocompleteHandler(s, i)
				return
			}
		}
	})
	return nil
}
