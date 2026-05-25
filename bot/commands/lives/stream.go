package lives

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// StartMonitor inicia o loop de checagem em segundo plano
func StartMonitor(s *discordgo.Session) {
	// Cria um ticker que "bate" a cada 5 minutos
	ticker := time.NewTicker(5 * time.Minute)

	// Executa a goroutine em background
	go func() {
		for {
			select {
			case <-ticker.C:
				// Aqui você busca no MongoDB a lista de streamers cadastrados
				// E faz as requisições para as APIs da Twitch e YouTube
				fmt.Println("Checando status das lives...")
				checkTwitchLives(s)
				checkYouTubeLives(s)
			}
		}
	}()
}
