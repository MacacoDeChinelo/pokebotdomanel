package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func RegisterSlashCommands(s *discordgo.Session) { //, guildID string) {

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "batalhar",
			Description: "Batalhe contra outro treinador!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "oponente",
					Description: "Usuário que você quer desafiar",
					Required:    true,
				},
			},
		},
		{
			Name:        "placar",
			Description: "Veja o placar de batalhas!",
		},
		{
			Name:        "sortear",
			Description: "Sorteia Pokémons diários para o servidor",
		},
	}

	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			fmt.Println("Erro ao criar comando:", err)
		}
	}
}
