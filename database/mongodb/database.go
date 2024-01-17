package mongodb

import (
	"context"
	"fmt"
	"os"
	"pidorator-bot/database"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	c   *mongo.Client
	log *zerolog.Logger
}

func New(ctx context.Context, log *zerolog.Logger) Database {

	result := Database{log: log}

	err := result.NewConnection(ctx)
	if err != nil {
		log.Error().Msgf("Error. Can't connect to database. err = %s", err)
	}

	return result
}

func (DB *Database) NewConnection(ctx context.Context) error {
	adress := os.Getenv("DBADRESS")
	username := os.Getenv("DBUSERNAME")
	password := os.Getenv("DBPASSWORD")

	// format: mongodb://login:password@adress
	uri := fmt.Sprintf("mongodb://%s:%s@%s", username, password, adress)
	clientOptions := options.Client().ApplyURI(uri)
	nameOpt := options.Client().SetAppName("Pidorator-bot")

	// TODO check which one is more efficient/secure
	timeoutOpt := options.Client().SetTimeout(25 * time.Second)
	connCtx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	client, err := mongo.Connect(connCtx, clientOptions, nameOpt, timeoutOpt)
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(connCtx, nil)
	if err != nil {
		return err
	}

	DB.c = client

	DB.log.Info().Msg("Connected to MongoDB!")

	return nil
}

func (DB *Database) Disconnect(ctx context.Context) {
	if DB.c == nil {
		DB.log.Error().Msgf("Error. Database client not found")
		return
	}
	err := DB.c.Disconnect(ctx)
	if err != nil {
		DB.log.Error().Msgf("Error. Can't disconnect from database. err = %s", err)
		return
	}

	DB.log.Debug().Msg("Connection to MongoDB is also closed.")
}

func (DB *Database) GetAllPlayers(ctx context.Context) ([]*database.GameData, error) {
	c := DB.c.Database("discord-bot").Collection("pidorator-game")

	findOptions := options.Find()

	var results []*database.GameData

	cur, err := c.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		DB.log.Error().Err(err)
		return []*database.GameData{}, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem database.GameData
		err := cur.Decode(&elem)
		if err != nil {
			DB.log.Error().Err(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		DB.log.Error().Err(err)
		return []*database.GameData{}, err
	}

	DB.log.Debug().Msgf("Found multiple documents (array of pointers): %+v\n", results)

	return results, nil
}

func (DB *Database) IncreaseUserScore(ctx context.Context, guildID string, userID string) error {
	c := DB.c.Database("discord-bot").Collection("pidorator-game")
	filter := bson.D{{Key: "userID", Value: userID}, {Key: "guildID", Value: guildID}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "score", Value: 1}}}}

	result, err := c.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	DB.log.Debug().Msgf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)

	return nil
}
