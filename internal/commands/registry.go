package commands

import "github.com/bwmarrin/discordgo"

// Command define a interface para todos os comandos
type Command struct {
	Definition *discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// CommandRegistry armazena todos os comandos registrados
var CommandRegistry = make(map[string]Command)

// RegisterCommand adiciona um comando ao registro
func RegisterCommand(cmd Command) {
	CommandRegistry[cmd.Definition.Name] = cmd
}
