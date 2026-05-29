package utils

import (
	"log"

	libdatabase "github.com/MacacoDeChinelo/readconf/pkg/database"
	"github.com/bwmarrin/discordgo"
)

func Init() (session *discordgo.Session) {

	err := libdatabase.Connect(0)
	if err != nil {
		log.Fatal("Erro ao conectar na base de dados:", err)
		return
	}
	err = libdatabase.ConnectDBConfig("BotPokemon")
	if err != nil {
		log.Fatal("Erro ao conectar na base BotPokemon:", err)
		return
	}

	tokenEncript := libdatabase.GetVariable("tokenDiscord").(string)

	token, err := Decrypt(tokenEncript)
	if err != nil {
		log.Fatal("Erro ao descriptografar o token do Discord:", err)
		return
	}

	session, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Erro ao criar sessão do Discord:", err)
	}

	return session
}

//https://github.com/MacacoDeChinelo/readconf
