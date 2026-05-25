package commands

import (
	"fmt"
	"math/rand"
	"time"

	"pokebot/database"
	"pokebot/models"

	"github.com/bwmarrin/discordgo"
)

func HandleSortear(s *discordgo.Session, m *discordgo.MessageCreate) {
	hoje := time.Now().Format("2006-01-02")

	// Verifica se o usuário que chamou já tem um sorteio hoje
	userScore, err := database.GetDailyScore(m.GuildID, m.Author.ID, hoje)

	if err != nil { // Se não tem, faz o sorteio geral
		members, err := s.GuildMembers(m.GuildID, "", 1000) // Requer intent de GuildMembers no Discord Developer Portal
		if err == nil {
			for _, member := range members {
				if member.User.Bot {
					continue
				}

				// Verifica se esse membro já tem
				_, errCheck := database.GetDailyScore(m.GuildID, member.User.ID, hoje)
				if errCheck != nil { // Sorteia para o membro
					pool, errPool := database.GetRandomPokemon()
					if errPool != nil {
						continue
					}

					power := 0
					if pool.Raridade == "normal" {
						power = rand.Intn(30) + 1 // 1 a 30
					} else { // legendary
						power = rand.Intn(21) + 20 // 20 a 40
					}

					newScore := &models.PokemonScore{
						DataSorteio: hoje,
						ServerID:    m.GuildID,
						UserID:      member.User.ID,
						Pokemon:     pool.Nome,
						Raridade:    pool.Raridade,
						Power:       power,
						Vitorias:    0,
						Derrotas:    0,
						URL:         pool.URL,
					}
					database.SaveDailyScore(newScore)

					// Salva a referência do chamador para exibir depois
					if member.User.ID == m.Author.ID {
						userScore = newScore
					}
				}
			}
		}
	}

	if userScore == nil { // Fallback de segurança se falhar ao buscar o próprio caller na lista
		s.ChannelMessageSend(m.ChannelID, "Ocorreu um erro ao buscar seu Pokémon.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎉 Seu Pokémon Diário!",
		Description: fmt.Sprintf("<@%s>, você tirou um **%s**!", m.Author.ID, userScore.Pokemon),
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Raridade", Value: userScore.Raridade, Inline: true},
			{Name: "Power", Value: fmt.Sprintf("%d", userScore.Power), Inline: true},
			{Name: "Vitórias / Derrotas", Value: fmt.Sprintf("%d / %d", userScore.Vitorias, userScore.Derrotas), Inline: false},
		},
		Image: &discordgo.MessageEmbedImage{URL: userScore.URL},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
