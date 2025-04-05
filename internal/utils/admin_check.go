package utils

import (
	"github.com/bwmarrin/discordgo"
)

// HasAdminPermission verifica se o usuário tem permissão de administrador
func HasAdminPermission(s *discordgo.Session, member *discordgo.Member, guildID string) bool {
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err == nil && (role.Permissions&discordgo.PermissionAdministrator) != 0 {
			return true
		}
	}
	return false
}
