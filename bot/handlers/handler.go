package handlers

import (
	"fmt"
	"strings"

	"pokebot/bot/commands"

	"github.com/bwmarrin/discordgo"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Fields(m.Content)
	if len(args) == 0 {
		return
	}

	command := strings.ToLower(args[0])

	switch command {
	case "!sortear":
		commands.HandleSortear(s, m)
	case "!batalhar":
		commands.HandleBatalhar(s, m, args)
	case "!placar":
		commands.HandlePlacar(s, m)
	}
}

// Evento disparado quando o bot entra em um servidor
func OnGuildCreateHandler(s *discordgo.Session, g *discordgo.GuildCreate) {

	if isAllowed(g.Guild.ID) {
		return
	}

	fmt.Printf(
		"Servidor não autorizado detectado: %s (%s)\n",
		g.Guild.Name,
		g.Guild.ID,
	)

	err := s.GuildLeave(g.Guild.ID)
	if err != nil {
		fmt.Println("Erro ao sair do servidor:", err)
	}
}

func OnSlashCommandCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	switch data.Name {

	case "batalhar":
		commands.HandleBatalharSlash(s, i)
	case "placar":
		commands.HandlePlacarSlash(s, i)
	case "sortear":
		commands.HandleSortearSlash(s, i)
	}
}
