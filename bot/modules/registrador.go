package commands

import (
	"fmt"
	pokemon "pokebot/bot/modules/pokemon"
	"pokebot/bot/modules/youtube"

	"github.com/bwmarrin/discordgo"
)

func RegisterSlashCommands(s *discordgo.Session) {
	// Agrupa todos os módulos em um slice de slices
	commandGroups := [][]*discordgo.ApplicationCommand{
		pokemon.GetCommands(),
		youtube.GetCommands(),
		// Adicione novos módulos aqui no futuro (ex: economia.GetCommands(), musica.GetCommands())
	}

	// Percorre os grupos e registra cada comando
	for _, group := range commandGroups {
		for _, cmd := range group {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
			if err != nil {
				fmt.Printf("Erro ao criar comando '%s': %v\n", cmd.Name, err)
			} else {
				fmt.Printf("Comando '%s' registrado com sucesso!\n", cmd.Name)
			}
		}
	}
}
