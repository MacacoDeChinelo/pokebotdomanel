package youtube

// bot/youtube/commands.go

import "github.com/bwmarrin/discordgo"

// Retorna os comandos deste módulo para o main.go
func GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "setyoutube",
			Description: "Configura o alerta de live do YouTube para este servidor",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "canal_texto",
					Description: "Canal do Discord onde o aviso será enviado",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "cargo",
					Description: "Cargo que será mencionado (Ex: @Inscritos)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id_youtube",
					Description: "ID do canal do YouTube (Ex: UCxxxxxxxxxxxx)",
					Required:    true,
				},
			},
		},
	}
}
