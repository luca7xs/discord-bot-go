package config

import "fmt"

const (
	// Bot
	Token = ""

	// Database
	User     = "root"
	Password = "password"
	Host     = "localhost"
	Port     = "3306"
	Database = "discord_bot"
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
