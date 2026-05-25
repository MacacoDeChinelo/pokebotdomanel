package utils

import (
	"fmt"
	"log"
	"os"
	"pokebot/database"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func Init() (config Config, session *discordgo.Session) {
	config = Config{}
	// Carrega o .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}
	mongoURI := os.Getenv("MONGO_URI")
	config.MongoURI = mongoURI

	//fmt.Println(mongoURI)
	if err := database.Connect(mongoURI); err != nil {
		log.Fatal("Erro ao conectar no MongoDB:", err)
	}
	fmt.Println("Conectado ao MongoDB!")

	token := os.Getenv("DISCORD_TOKEN")
	config.Token = token
	//fmt.Println(token)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Erro ao criar sessão do Discord:", err)
	}

	return config, dg
}
