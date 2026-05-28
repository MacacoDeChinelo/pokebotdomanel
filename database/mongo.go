package database

import (
	"context"
	"log"
	"pokebot/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var DB *mongo.Database

func InitDB(db *mongo.Database) {
	DB = db

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

func UpdateBattleResult(winnerID, loserID bson.ObjectID) {
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

func InsertYouTubeAlert(ctx context.Context, alerta models.YouTubeAlert) error {
	// Acessando a collection streamer_alerts dentro do banco BotPokemon [cite: 8, 76]
	// Substitua 'MongoClient' pela sua variável real de conexão com o MongoDB
	collection := DB.Collection("streamer_alerts")

	_, err := collection.InsertOne(ctx, alerta)
	if err != nil {
		log.Printf("Erro ao inserir configuração de alerta do YouTube: %v", err)
		return err
	}

	return nil
}

func GetAllYouTubeAlerts(ctx context.Context) ([]models.YouTubeAlert, error) {
	var alertas []models.YouTubeAlert

	collection := DB.Collection("streamer_alerts")

	// Usamos bson.M{} vazio para buscar todos os documentos da collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Erro ao buscar alertas do YouTube no banco: %v", err)
		return nil, err
	}

	// É boa prática garantir que o cursor seja fechado ao final da operação
	defer cursor.Close(ctx)

	// O método All itera sobre o cursor e decodifica tudo direto no slice 'alertas'
	if err = cursor.All(ctx, &alertas); err != nil {
		log.Printf("Erro ao decodificar os alertas do YouTube: %v", err)
		return nil, err
	}

	return alertas, nil
}

// UpdateLiveStatus altera a flag is_live de um alerta específico no MongoDB
func UpdateLiveStatus(ctx context.Context, id bson.ObjectID, isLive bool) error {
	// Acessando a collection streamer_alerts dentro do banco BotPokemon
	collection := DB.Collection("streamer_alerts")

	// Filtro: busca o documento exatamente pelo ID (ObjectID gerado pelo Mongo)
	filter := bson.M{"_id": id}

	// Operação: utiliza o operador $set para alterar apenas o campo is_live
	update := bson.M{
		"$set": bson.M{
			"is_live": isLive,
		},
	}

	// Executa a atualização de um único documento
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Erro ao atualizar o status de live no banco de dados para o ID %s: %v", id.Hex(), err)
		return err
	}

	return nil
}
