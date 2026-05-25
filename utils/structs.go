package utils

import "github.com/bwmarrin/discordgo"

type Config struct {
	Token        string
	MongoURI     string
	DatabaseName string
	Session      *discordgo.Session
}
