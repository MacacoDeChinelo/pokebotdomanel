package commands

import (
	pokemon "darthverde/bot/modules/pokemon"
	"darthverde/bot/modules/twitch"
	"darthverde/bot/modules/youtube"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func RegisterSlashCommands(s *discordgo.Session) {

	// 1. Cria um slice (array) vazio que vai receber TODOS os comandos
	var allCommands []*discordgo.ApplicationCommand
	// Agrupa todos os módulos em um slice de slices
	commandGroups := [][]*discordgo.ApplicationCommand{
		pokemon.GetCommands(),
		youtube.GetCommands(),
		twitch.GetCommands(),
		// Adicione novos módulos aqui no futuro (ex: economia.GetCommands(), musica.GetCommands())
	}
	// 3. Junta tudo dentro do nosso array único 'allCommands'
	for _, group := range commandGroups {
		allCommands = append(allCommands, group...)
	}
	// Percorre os grupos e registra cada comando
	// 4. Envia TODOS os comandos de uma só vez (Bulk Overwrite)
	fmt.Printf("Enviando lote com %d comandos para o Discord...\n", len(allCommands))
	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", allCommands)

	if err != nil {
		fmt.Printf("❌ Erro ao registrar comandos em lote: %v\n", err)
	} else {
		fmt.Println("✅ Todos os comandos foram registrados instantaneamente com sucesso!")
	}
}
