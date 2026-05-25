package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Verifica se o servidor é permitido
func isAllowed(guildID string) bool {
	_, ok := AllowedGuilds[guildID]
	return ok
}

// Lista de servidores permitidos
var AllowedGuilds = map[string]struct{}{
	"1508251298737295471": {}, //server do manel
	"1214208393661648917": {}, // server do bruno
	"555555555555555555":  {},
}

// Remove o bot de servidores não autorizados ao iniciar
func CheckGuilds(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		if !isAllowed(guild.ID) {

			fmt.Printf(
				"Saindo do servidor não autorizado: %s (%s)\n",
				guild.Name,
				guild.ID,
			)

			err := s.GuildLeave(guild.ID)
			if err != nil {
				fmt.Println("Erro ao sair:", err)
			}
		}
	}
}
