package commands

import (
	"fmt"
	"sort"
	"time"

	"pokebot/database"

	"github.com/bwmarrin/discordgo"
)

func HandlePlacar(s *discordgo.Session, m *discordgo.MessageCreate) {
	hoje := time.Now().Format("2006-01-02")
	scores, err := database.GetServerDailyLeaderboard(m.GuildID, hoje)
	if err != nil || len(scores) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Ninguém sorteou Pokémons hoje ainda.")
		return
	}

	// Ordena: vitórias menos derrotas descrescente
	sort.Slice(scores, func(i, j int) bool {
		scoreI := scores[i].Vitorias - scores[i].Derrotas
		scoreJ := scores[j].Vitorias - scores[j].Derrotas
		return scoreI > scoreJ
	})

	var lines []string
	for i, sc := range scores {
		medal := "🏅"
		if i == 0 {
			medal = "🥇"
		} else if i == 1 {
			medal = "🥈"
		} else if i == 2 {
			medal = "🥉"
		}

		//saldo := sc.Vitorias - sc.Derrotas
		linha := fmt.Sprintf("%s **<@%s>** - %s (Power: %d | Vitórias: %d | Derrotas: %d)", // | Saldo: %d)",
			medal, sc.UserID, sc.Pokemon, sc.Power, sc.Vitorias, sc.Derrotas) //, saldo)
		lines = append(lines, linha)
	}

	// Paginação para evitar limite de 4096 caracteres no Embed
	chunkSize := 15
	var chunks [][]string
	for i := 0; i < len(lines); i += chunkSize {
		end := i + chunkSize
		if end > len(lines) {
			end = len(lines)
		}
		chunks = append(chunks, lines[i:end])
	}

	for i, chunk := range chunks {
		desc := ""
		for _, line := range chunk {
			desc += line + "\n"
		}

		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("🏆 Placar Pokémon (Parte %d/%d)", i+1, len(chunks)),
			Description: desc,
			Color:       0xffd700,
		}

		// Adiciona a imagem do podium no último embed
		if i == len(chunks)-1 {
			embed.Image = &discordgo.MessageEmbedImage{
				URL: "https://png.pngtree.com/png-clipart/20230124/ourmid/pngtree-winner-podium-cartoon-png-image_6564495.png",
			}
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}

func HandlePlacarSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {

	hoje := time.Now().Format("2006-01-02")

	scores, err := database.GetServerDailyLeaderboard(i.GuildID, hoje)
	if err != nil || len(scores) == 0 {

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Ninguém sorteou Pokémons hoje ainda.",
			},
		})
		return
	}

	// ordenação
	sort.Slice(scores, func(a, b int) bool {
		return (scores[a].Vitorias - scores[a].Derrotas) >
			(scores[b].Vitorias - scores[b].Derrotas)
	})

	var lines []string

	for idx, sc := range scores {

		medal := "🏅"
		if idx == 0 {
			medal = "🥇"
		} else if idx == 1 {
			medal = "🥈"
		} else if idx == 2 {
			medal = "🥉"
		}

		linha := fmt.Sprintf(
			"%s **<@%s>** - %s (Power: %d | Vitórias: %d | Derrotas: %d)",
			medal,
			sc.UserID,
			sc.Pokemon,
			sc.Power,
			sc.Vitorias,
			sc.Derrotas,
		)

		lines = append(lines, linha)
	}

	// paginação
	chunkSize := 15
	var chunks [][]string

	for i := 0; i < len(lines); i += chunkSize {
		end := i + chunkSize
		if end > len(lines) {
			end = len(lines)
		}
		chunks = append(chunks, lines[i:end])
	}

	// responder com "loading" primeiro (boa prática)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// enviar embeds
	for iChunk, chunk := range chunks {

		desc := ""
		for _, line := range chunk {
			desc += line + "\n"
		}

		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("🏆 Placar Pokémon (Parte %d/%d)", iChunk+1, len(chunks)),
			Description: desc,
			Color:       0xffd700,
		}

		if iChunk == len(chunks)-1 {
			embed.Image = &discordgo.MessageEmbedImage{
				URL: "https://png.pngtree.com/png-clipart/20230124/ourmid/pngtree-winner-podium-cartoon-png-image_6564495.png",
			}
		}

		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{embed},
		})
	}
}
