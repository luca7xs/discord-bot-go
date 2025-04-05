package components

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// ComponentHandler define a interface para handlers de componentes
type ComponentHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

// Component define um componente registrado
type Component struct {
	CustomID string
	Handler  ComponentHandler
}

// ComponentRegistry armazena todos os componentes registrados
var ComponentRegistry = make(map[string]Component)

// RegisterComponent adiciona um componente ao registro
func RegisterComponent(customID string, handler ComponentHandler) error {
	if _, exists := ComponentRegistry[customID]; exists {
		return errors.New("CustomID já registrado: " + customID)
	}
	ComponentRegistry[customID] = Component{
		CustomID: customID,
		Handler:  handler,
	}
	return nil
}

// GetComponentHandler retorna o handler de um componente pelo CustomID
func GetComponentHandler(customID string) (ComponentHandler, bool) {
	if component, exists := ComponentRegistry[customID]; exists {
		return component.Handler, true
	}
	return nil, false
}
