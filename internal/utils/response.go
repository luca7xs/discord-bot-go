package utils

import (
	"github.com/bwmarrin/discordgo"
)

// ResponseOptions define as opções para a resposta em embed
type ResponseOptions struct {
	Title       string                         // Opcional: título do embed
	Description string                         // Opcional: descrição do embed
	Fields      []*discordgo.MessageEmbedField // Opcional: campos do embed
	Image       *discordgo.MessageEmbedImage   // Opcional: imagem do embed
	Color       int                            // Opcional: cor do embed
}

// RespondEphemeralEmbed cria uma resposta efêmera com um embed
func RespondEphemeralEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, opts ResponseOptions) error {
	embed := &discordgo.MessageEmbed{
		Title:       opts.Title,
		Description: opts.Description,
		Fields:      opts.Fields,
		Image:       opts.Image,
		Color:       opts.Color, // Cor padrão #2b2d31
	}

	// Se não for fornecido, definir a cor padrão
	if embed.Color == 0 {
		embed.Color = 0x2b2d31 // Usando a cor padrão da utils
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
