package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pokebot/bot"
	"pokebot/database"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {

	// Carrega o .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}
	mongoURI := os.Getenv("MONGO_URI")
	fmt.Println(mongoURI)
	if err := database.Connect(mongoURI); err != nil {
		log.Fatal("Erro ao conectar no MongoDB:", err)
	}
	fmt.Println("Conectado ao MongoDB!")

	token := os.Getenv("DISCORD_TOKEN")
	fmt.Println(token)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Erro ao criar sessão do Discord:", err)
	}

	// Intents necessários, incluindo GuildMembers para o sorteio varrer o server
	dg.Identify.Intents = discordgo.IntentsAll //discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	dg.AddHandler(bot.MessageHandler)

	if err = dg.Open(); err != nil {
		log.Fatal("Erro ao abrir conexão com o Discord:", err)
	}

	fmt.Println("Bot Pokémon está rodando! Pressione CTRL-C para sair.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
