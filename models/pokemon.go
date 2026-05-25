package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PokemonPool struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Nome     string             `bson:"nome"`
	Tipos    []string           `bson:"tipos"`
	Ataques  []AtaquesPool      `bson:"ataques"`
	Raridade string             `bson:"raridade"`
	URL      string             `bson:"url"`
}

type AtaquesPool struct {
	Nome string `bson:"nome"`
	Tipo string `bson:"tipo"`
}

type PokemonScore struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	DataSorteio string             `bson:"data_sorteio"`
	ServerID    string             `bson:"server_id"`
	UserID      string             `bson:"user_id"`
	Pokemon     string             `bson:"pokemon"`
	Raridade    string             `bson:"raridade"`
	Power       int                `bson:"power"`
	Vitorias    int                `bson:"vitorias"`
	Derrotas    int                `bson:"derrotas"`
	URL         string             `bson:"url"`
}
