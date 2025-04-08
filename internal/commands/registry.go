package commands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// CommandHandler define a interface para handlers de comandos
type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

// Command define a interface para todos os comandos
type Command struct {
	Definition          *discordgo.ApplicationCommand
	Handler             CommandHandler
	AutocompleteHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) // <- Suporte adicionado
}

// CommandRegistry armazena todos os comandos registrados
var CommandRegistry = make(map[string]Command)

// RegisterCommand adiciona um comando ao registro
func RegisterCommand(cmd Command) error {
	// Verifica se o comando já existe
	if _, exists := CommandRegistry[cmd.Definition.Name]; exists {
		return errors.New("Comando já registrado: " + cmd.Definition.Name)
	}
	CommandRegistry[cmd.Definition.Name] = cmd
	return nil
}
