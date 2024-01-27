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

const (
	BotDBName          string = "discord-bot"
	GameCollectionName string = "pidorator-game"
	ErrorNotExist      string = "mongo: no documents in result"
)

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

func (DB *Database) AddPlayer(ctx context.Context, data *database.PlayerData) error {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return err
	}

	p, err := DB.FindPlayer(ctx, data.GuildID, data.UserID)
	if err != nil && err.Error() != ErrorNotExist {
		return err
	}

	if p.GuildID != "" {
		err := fmt.Errorf(database.PlayerAlreadyExistError)
		return err
	}

	c := DB.c.Database(BotDBName).Collection(GameCollectionName)

	result, err := c.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}

	if result == nil {
		err := fmt.Errorf("error. Player not added")
		return err
	}

	return nil
}

func (DB *Database) GetAllPlayers(ctx context.Context, guildID string) ([]*database.PlayerData, error) {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return []*database.PlayerData{}, err
	}

	c := DB.c.Database(BotDBName).Collection(GameCollectionName)

	findOptions := options.Find().SetSort(bson.D{{Key: "score", Value: -1}})

	var results []*database.PlayerData

	cur, err := c.Find(ctx, bson.D{{Key: "guildID", Value: guildID}}, findOptions)
	if err != nil {
		return []*database.PlayerData{}, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem database.PlayerData
		err := cur.Decode(&elem)
		if err != nil {
			DB.log.Error().Err(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return []*database.PlayerData{}, err
	}

	DB.log.Debug().Msgf("Found multiple documents(%d). Array of pointers: %+v", len(results), results)

	return results, nil
}

func (DB *Database) FindPlayer(ctx context.Context, guildID string, userID string) (*database.PlayerData, error) {
	result := database.PlayerData{}
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return &result, err
	}

	c := DB.c.Database(BotDBName).Collection(GameCollectionName)

	findOptions := options.FindOne()

	filter := bson.D{{Key: "guildID", Value: guildID}, {Key: "userID", Value: userID}}

	err := c.FindOne(ctx, filter, findOptions).Decode(&result)
	if err != nil {
		return &database.PlayerData{}, err
	}

	return &result, nil
}

func (DB *Database) IncreasePlayerScore(ctx context.Context, guildID string, userID string) error {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return err
	}

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

func (DB *Database) UpdatePlayersData(ctx context.Context, guildID string, data []*database.PlayerData) error {

	return nil
}
