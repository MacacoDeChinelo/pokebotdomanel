package lives

import "github.com/bwmarrin/discordgo"

func GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "addstream",
			Description: "Adiciona um canal da Twitch ou YouTube para ser monitorado",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "plataforma",
					Description: "twitch ou youtube",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "canal",
					Description: "Nome do streamer ou ID do canal",
					Required:    true,
				},
			},
		},
		{
			Name:        "removestream",
			Description: "Remove um canal do monitoramento",
		},
	}
}
