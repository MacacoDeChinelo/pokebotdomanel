package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Estrutura baseada no seu modelo do YouTube
type TwitchAlert struct {
	ID             bson.ObjectID `bson:"_id,omitempty"`
	ServerID       string        `bson:"server_id"`
	DiscordChannel string        `bson:"discord_channel"`
	RoleToMention  string        `bson:"role_to_mention"`
	TwitchChannel  string        `bson:"twitch_channel"`
	IsLive         bool          `bson:"is_live"`
}
