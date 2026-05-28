// models/youtube.go
package models

import "go.mongodb.org/mongo-driver/v2/bson"

type YouTubeAlert struct {
	ID             bson.ObjectID `bson:"_id,omitempty"`
	ServerID       string        `bson:"server_id"`
	DiscordChannel string        `bson:"discord_channel"`
	RoleToMention  string        `bson:"role_to_mention"`
	YouTubeChannel string        `bson:"youtube_channel"`
	IsLive         bool          `bson:"is_live"`
}
