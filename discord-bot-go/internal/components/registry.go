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
func RegisterComponent(component Component) error {
	// Verifica se o componente já existe
	if _, exists := ComponentRegistry[component.CustomID]; exists {
		return errors.New("CustomID já registrado: " + component.CustomID)
	}
	ComponentRegistry[component.CustomID] = component
	return nil
}
