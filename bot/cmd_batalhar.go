package bot

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"pokebot/database"

	"github.com/bwmarrin/discordgo"
)

var (
	cooldowns = make(map[string]time.Time)
	cdMutex   sync.Mutex
)

func HandleBatalhar(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(m.Mentions) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Mencione o treinador que deseja batalhar! Ex: `!batalhar @usuario`")
		return
	}
	target := m.Mentions[0]
	if target.ID == m.Author.ID || target.Bot {
		s.ChannelMessageSend(m.ChannelID, "Você não pode batalhar com você mesmo ou com bots!")
		return
	}

	// Verifica Cooldown (10 minutos)
	cdMutex.Lock()
	if lastUse, exists := cooldowns[m.Author.ID]; exists {
		if time.Since(lastUse) < 10*time.Minute {
			faltam := (10 * time.Minute) - time.Since(lastUse)
			cdMutex.Unlock()
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⏳ Seu Pokémon está no Centro Pokémon curando! Aguarde %d minutos.", int(faltam.Minutes())))
			return
		}
	}
	cooldowns[m.Author.ID] = time.Now()
	cdMutex.Unlock()

	hoje := time.Now().Format("2006-01-02")
	p1, err1 := database.GetDailyScore(m.GuildID, m.Author.ID, hoje)
	p2, err2 := database.GetDailyScore(m.GuildID, target.ID, hoje)

	if err1 != nil || err2 != nil {
		s.ChannelMessageSend(m.ChannelID, "Ambos os jogadores precisam sortear seus Pokémons hoje (`!sortear`) antes de batalhar.")
		return
	}

	// Lógica de Batalha (Loop interno)
	hp1, hp2 := p1.Power, p2.Power
	winner, loser := p1, p2

	for hp1 > 0 && hp2 > 0 {
		// P1 Ataca
		hp2 -= hit(p1.Raridade)
		if hp2 <= 0 {
			winner, loser = p1, p2
			break
		}
		// P2 Ataca
		hp1 -= hit(p2.Raridade)
		if hp1 <= 0 {
			winner, loser = p2, p1
			break
		}
	}

	database.UpdateBattleResult(winner.ID, loser.ID)

	embed := &discordgo.MessageEmbed{
		Title:       "⚔️ Resultado da Batalha!",
		Description: fmt.Sprintf("<@%s> (%s) VS <@%s> (%s)", p1.UserID, p1.Pokemon, p2.UserID, p2.Pokemon),
		Color:       0xff0000,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Vencedor!", Value: fmt.Sprintf("<@%s>\nO **%s** ganhou +1 de Power!", winner.UserID, winner.Pokemon), Inline: false},
		},
		Image: &discordgo.MessageEmbedImage{URL: winner.URL},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func hit(raridade string) int {
	if raridade == "legendary" {
		return rand.Intn(6) + 5 // 5 a 10
	}
	return rand.Intn(10) + 1 // 1 a 10
}
