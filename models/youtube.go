// models/youtube.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type YouTubeAlert struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	ServerID       string             `bson:"server_id"`
	DiscordChannel string             `bson:"discord_channel"` // Onde o bot vai mandar a mensagem
	RoleToMention  string             `bson:"role_to_mention"` // Cargo para pingar (@everyone ou um cargo específico)
	YouTubeChannel string             `bson:"youtube_channel"` // ID do canal do YouTube
	IsLive         bool               `bson:"is_live"`         // Controle para evitar flood de mensagens
}
