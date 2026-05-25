package utils

import (
	"fmt"
	"log"
	"os"
	"pokebot/database"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func Init() (session *discordgo.Session) {

	// Carrega o .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}
	mongoURI := os.Getenv("MONGO_URI")

	//fmt.Println(mongoURI)
	if err := database.Connect(mongoURI); err != nil {
		log.Fatal("Erro ao conectar na base de dados:", err)
	}
	fmt.Println("Conectado à base de dados!")

	token := os.Getenv("DISCORD_TOKEN")

	//fmt.Println(token)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Erro ao criar sessão do Discord:", err)
	}

	return dg
}
