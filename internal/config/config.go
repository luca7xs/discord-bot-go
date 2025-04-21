package config

import "fmt"

const (
	// Bot
	Token   = "MTM1NzUwNDc4ODU4OTQ0NTI4MQ.GvwNsg.l9YsHuXV1kdHvmeHXJ9lCKS4c0bdppO8BBLLlY"
	GuildID = "1357505318115999935"

	// Database
	User     = "root"
	Password = ""
	Host     = "192.168.0.2"
	Port     = "3306"
	Database = "discord_bot_go"
)

var DatabaseURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", User, Password, Host, Port, Database)

var TicketTypes = []struct {
	Label       string
	Value       string
	Description string
}{
	{
		Label:       "Suporte",
		Value:       "support",
		Description: "Para questões de suporte geral",
	},
	{
		Label:       "Denúncia",
		Value:       "report",
		Description: "Para denunciar um problema ou usuário",
	},
	{
		Label:       "Sugestão",
		Value:       "suggestion",
		Description: "Para sugerir algo ao servidor",
	},
}
