package utils

import (
	"github.com/bwmarrin/discordgo"
)

// FetchAllMessages busca todas as mensagens de um canal com paginação
func FetchAllMessages(s *discordgo.Session, channelID string) ([]*discordgo.Message, error) {
	var allMessages []*discordgo.Message
	lastID := ""

	for {
		messages, err := s.ChannelMessages(channelID, 100, lastID, "", "")
		if err != nil {
			return nil, err
		}
		if len(messages) == 0 {
			break
		}
		allMessages = append(allMessages, messages...)
		lastID = messages[len(messages)-1].ID
	}

	return allMessages, nil
}
