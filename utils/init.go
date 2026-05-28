package utils

import (
	"log"
	"pokebot/database"

	"github.com/bwmarrin/discordgo"
	libdatabase "github.com/jolealpe89/readconf/pkg/database"
)

// func Init() (session *discordgo.Session) {
//
//		// Carrega o .env
//		err := godotenv.Load()
//		if err != nil {
//			log.Fatal("Erro ao carregar .env")
//		}
//		mongoURI := os.Getenv("MONGO_URI")
//
//		//fmt.Println(mongoURI)
//		if err := database.Connect(mongoURI); err != nil {
//			log.Fatal("Erro ao conectar na base de dados:", err)
//		}
//		fmt.Println("Conectado à base de dados!")
//
//		token := os.Getenv("DISCORD_TOKEN")
//
//		//fmt.Println(token)
//		dg, err := discordgo.New("Bot " + token)
//		if err != nil {
//			log.Fatal("Erro ao criar sessão do Discord:", err)
//		}
//
//		return dg
//	}

func Init() (session *discordgo.Session) {

	// Carrega o .env
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Erro ao carregar .env")
	//}

	//db, err := libdatabase.ConnectDB("BotPokemon")
	err := libdatabase.Connect()
	if err != nil {
		log.Fatal("Erro ao conectar na base de dados:", err)
		return
	}
	err = libdatabase.ConnectDBConfig("BotPokemon")
	if err != nil {
		log.Fatal("Erro ao conectar na base BotPokemon:", err)
		return
	}
	database.InitDB(libdatabase.DB)
	tokenEncript := libdatabase.GetVariable("tokenDiscord").(string)

	token, err := Decrypt(tokenEncript)
	if err != nil {
		log.Fatal("Erro ao descriptografar o token do Discord:", err)
		return
	}

	//token := os.Getenv("DISCORD_TOKEN")

	//fmt.Println(token)
	session, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Erro ao criar sessão do Discord:", err)
	}

	return session
}

//https://github.com/jolealpe89/readconf
