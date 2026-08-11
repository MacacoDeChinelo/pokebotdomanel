package database

import (
	"context"
	"darthverde/models"

	libdatabase "github.com/MacacoDeChinelo/readconf/pkg/database"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Salva um novo alerta da Twitch
func InsertTwitchAlert(ctx context.Context, alerta models.TwitchAlert) error {
	collection := libdatabase.DB.Collection("twitch_alerts")
	_, err := collection.InsertOne(ctx, alerta)
	return err
}

// Pega todos os alertas para o monitor verificar
func GetAllTwitchAlerts(ctx context.Context) ([]models.TwitchAlert, error) {
	var alertas []models.TwitchAlert
	collection := libdatabase.DB.Collection("twitch_alerts")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &alertas); err != nil {
		return nil, err
	}
	return alertas, nil
}

// Atualiza se a live está ON ou OFF
func UpdateTwitchLiveStatus(ctx context.Context, id bson.ObjectID, isLive bool) error {
	collection := libdatabase.DB.Collection("twitch_alerts")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"is_live": isLive}}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
