package database

import (
	"context"
	"pokebot/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	DB = client.Database("BotPokemon")
	return nil
}

func GetRandomPokemon() (*models.PokemonPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}}}
	cursor, err := DB.Collection("pokemon_pool").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pokemons []models.PokemonPool
	if err = cursor.All(ctx, &pokemons); err != nil || len(pokemons) == 0 {
		return nil, err
	}
	return &pokemons[0], nil
}

func GetPokemon(name string) (*models.PokemonPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"nome": name}
	var pokemon models.PokemonPool
	err := DB.Collection("pokemon_pool").FindOne(ctx, filter).Decode(&pokemon)
	return &pokemon, err
}

func GetDailyScore(serverID, userID, date string) (*models.PokemonScore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var score models.PokemonScore
	filter := bson.M{"server_id": serverID, "user_id": userID, "data_sorteio": date}
	err := DB.Collection("pokemon_scores").FindOne(ctx, filter).Decode(&score)
	return &score, err
}

func SaveDailyScore(score *models.PokemonScore) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := DB.Collection("pokemon_scores").InsertOne(ctx, score)
	return err
}

func UpdateBattleResult(winnerID, loserID primitive.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Vencedor: +1 vitória, +1 power
	DB.Collection("pokemon_scores").UpdateOne(ctx,
		bson.M{"_id": winnerID},
		bson.M{"$inc": bson.M{"vitorias": 1, "power": 1}},
	)
	// Perdedor: +1 derrota
	DB.Collection("pokemon_scores").UpdateOne(ctx,
		bson.M{"_id": loserID},
		bson.M{"$inc": bson.M{"derrotas": 1}},
	)
}

func GetServerDailyLeaderboard(serverID, date string) ([]models.PokemonScore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"server_id": serverID, "data_sorteio": date}
	cursor, err := DB.Collection("pokemon_scores").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scores []models.PokemonScore
	if err = cursor.All(ctx, &scores); err != nil {
		return nil, err
	}
	return scores, nil
}
