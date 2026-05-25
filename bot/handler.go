package bot

import (
	"strings"

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
		HandleSortear(s, m)
	case "!batalhar":
		HandleBatalhar(s, m, args)
	case "!placar":
		HandlePlacar(s, m)
	}
}
