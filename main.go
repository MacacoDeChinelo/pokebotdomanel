package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	bot "darthverde/bot/handlers"
	commands "darthverde/bot/modules"
	"darthverde/bot/modules/twitch"
	"darthverde/bot/modules/youtube"
	"darthverde/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {
	var err error

	dg := utils.Init() // função que inicializa o banco de dados e retorna a sessão do Discord

	// Intents necessários, incluindo GuildMembers para o sorteio varrer o server
	dg.Identify.Intents = discordgo.IntentsAll //discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	dg.AddHandler(bot.MessageHandler)

	dg.AddHandler(bot.OnGuildCreateHandler)

	dg.AddHandler(bot.OnSlashCommandCreateHandler)

	if err = dg.Open(); err != nil {
		log.Fatal("Erro ao abrir conexão com o Discord:", err)
		return
	}
	botName := dg.State.User.Username

	// Registra os comandos de barra
	commands.RegisterSlashCommands(dg) // Registra os comandos de barra
	// Verifica os servidores atuais
	bot.CheckGuilds(dg)
	youtube.StartYouTubeMonitor(dg)
	twitch.StartTwitchMonitor(dg)
	fmt.Println(botName + " está pronto! Pressione CTRL-C para sair.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
