package commands

import (
	"discord-bot-go/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand(Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "avatar",
			Description: "Mostra o avatar de um usuário",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "usuario",
					Description: "Usuário para ver o avatar",
					Required:    false,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.Member.User
			if len(i.ApplicationCommandData().Options) > 0 {
				user = i.ApplicationCommandData().Options[0].UserValue(s)
			}
			avatarURL := user.AvatarURL("256")

			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Description: "Aqui está o avatar!",
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Usuário", Value: user.Username, Inline: true},
					{Name: "ID", Value: user.ID, Inline: true},
				},
				Image: &discordgo.MessageEmbedImage{
					URL: avatarURL,
				},
			})
		},
	})
}
