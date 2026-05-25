package commands

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

func HandlePokebattle(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(m.Mentions) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Mencione o treinador que deseja batalhar! Ex: `!pokebattle @usuario`")
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

// normal 🟢 raro 🟨 e lendário 🔶
func hit(raridade string) int {
	if raridade == "legendary 🔶" {
		return rand.Intn(6) + 5 // 5 a 10
	}
	return rand.Intn(10) + 1 // 1 a 10
}

func HandlePokebattleSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// =========================
	// VALIDAÇÕES
	// =========================
	if i.Member == nil || i.Member.User == nil {
		return
	}

	options := i.ApplicationCommandData().Options

	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Escolha um oponente para batalhar.",
			},
		})
		return
	}

	target := options[0].UserValue(s)

	if target == nil {
		return
	}

	if target.ID == i.Member.User.ID || target.Bot {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Você não pode batalhar consigo mesmo ou contra bots.",
			},
		})
		return
	}

	// =========================
	// COOLDOWN
	// =========================
	cdMutex.Lock()

	if lastUse, exists := cooldowns[i.Member.User.ID]; exists {

		if time.Since(lastUse) < 10*time.Minute {

			faltam := (10 * time.Minute) - time.Since(lastUse)

			cdMutex.Unlock()

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						"⏳ Seu Pokémon está no Centro Pokémon curando.\nAguarde %d minutos.",
						int(faltam.Minutes()),
					),
				},
			})

			return
		}
	}

	cooldowns[i.Member.User.ID] = time.Now()

	cdMutex.Unlock()

	// =========================
	// BUSCA POKÉMONS
	// =========================
	hoje := time.Now().Format("2006-01-02")

	p1, err1 := database.GetDailyScore(
		i.GuildID,
		i.Member.User.ID,
		hoje,
	)

	p2, err2 := database.GetDailyScore(
		i.GuildID,
		target.ID,
		hoje,
	)

	if err1 != nil || err2 != nil || p1 == nil || p2 == nil {

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Ambos os jogadores precisam usar `/sortear` antes da batalha.",
			},
		})

		return
	}

	// =========================
	// RESPOSTA INICIAL
	// =========================
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// =========================
	// HP BASE
	// =========================
	hp1 := p1.Power * 10
	hp2 := p2.Power * 10

	winner := p1
	loser := p2

	// =========================
	// MENSAGEM INICIAL
	// =========================
	s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: fmt.Sprintf(
			"⚔️ **%s** VS **%s** ⚔️\n🔥 A batalha começou!",
			p1.Pokemon,
			p2.Pokemon,
		),
	})

	time.Sleep(2 * time.Second)

	// =========================
	// LOOP DE BATALHA
	// =========================
	turn := 1

	for hp1 > 0 && hp2 > 0 {

		// =====================
		// P1 ATACA
		// =====================
		damage1 := hit(p1.Raridade)

		critical1 := rand.Intn(100) < 20

		if critical1 {
			damage1 *= 2
		}

		hp2 -= damage1

		if hp2 < 0 {
			hp2 = 0
		}

		msg1 := fmt.Sprintf(
			"🎯 **Turno %d**\n\n⚡ **%s** atacou **%s**!\n💥 Dano: **%d**",
			turn,
			p1.Pokemon,
			p2.Pokemon,
			damage1,
		)

		if critical1 {
			msg1 += "\n✨ **ATAQUE CRÍTICO!**"
		}

		msg1 += fmt.Sprintf(
			"\n❤️ HP %s: **%d**",
			p2.Pokemon,
			hp2,
		)

		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: msg1,
		})

		time.Sleep(2 * time.Second)

		if hp2 <= 0 {
			winner = p1
			loser = p2
			break
		}

		// =====================
		// P2 ATACA
		// =====================
		dodge := rand.Intn(100) < 15

		if dodge {

			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Content: fmt.Sprintf(
					"💨 **%s desviou do ataque de %s!**",
					p1.Pokemon,
					p2.Pokemon,
				),
			})

			time.Sleep(2 * time.Second)

		} else {

			damage2 := hit(p2.Raridade)

			critical2 := rand.Intn(100) < 20

			if critical2 {
				damage2 *= 2
			}

			hp1 -= damage2

			if hp1 < 0 {
				hp1 = 0
			}

			msg2 := fmt.Sprintf(
				"🎯 **Turno %d**\n\n🔥 **%s** atacou **%s**!\n💥 Dano: **%d**",
				turn,
				p2.Pokemon,
				p1.Pokemon,
				damage2,
			)

			if critical2 {
				msg2 += "\n✨ **ATAQUE CRÍTICO!**"
			}

			msg2 += fmt.Sprintf(
				"\n❤️ HP %s: **%d**",
				p1.Pokemon,
				hp1,
			)

			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Content: msg2,
			})

			time.Sleep(2 * time.Second)

			if hp1 <= 0 {
				winner = p2
				loser = p1
				break
			}
		}

		turn++
	}

	// =========================
	// RESULTADO FINAL
	// =========================
	database.UpdateBattleResult(winner.ID, loser.ID)

	embed := &discordgo.MessageEmbed{
		Title: "🏆 Resultado da Batalha!",
		Description: fmt.Sprintf(
			"⚔️ **%s** venceu a batalha contra **%s**!",
			winner.Pokemon,
			loser.Pokemon,
		),
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Treinador vencedor",
				Value:  fmt.Sprintf("<@%s>", winner.UserID),
				Inline: false,
			},
			{
				Name:   "Recompensa",
				Value:  "+1 Power",
				Inline: false,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: winner.URL,
		},
	}

	s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})
}
