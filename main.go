package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	bot "pokebot/bot/handlers"
	"pokebot/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {
	var err error

	_, dg := utils.Init() // função que inicializa o banco de dados e retorna a sessão do Discord

	// Intents necessários, incluindo GuildMembers para o sorteio varrer o server
	dg.Identify.Intents = discordgo.IntentsAll //discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	dg.AddHandler(bot.MessageHandler)

	dg.AddHandler(bot.OnGuildCreateHandler)

	if err = dg.Open(); err != nil {
		log.Fatal("Erro ao abrir conexão com o Discord:", err)
	}

	fmt.Println("Bot Pokémon está rodando! Pressione CTRL-C para sair.")
	// Verifica os servidores atuais
	bot.CheckGuilds(dg)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
