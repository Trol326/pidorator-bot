package mongodb

import (
	"context"
	"fmt"
	"os"
	"pidorator-bot/app/database/model"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	c       *mongo.Client
	log     *zerolog.Logger
	botData map[string]*model.BotData
}

const (
	BotDBName string = "discord-bot"

	DataCollectionName  string = "pidorator"
	GameCollectionName  string = "pidorator-game"
	EventCollectionName string = "pidorator-events"
	ErrorNotExist       string = "mongo: no documents in result"
)

func New(ctx context.Context, log *zerolog.Logger) Database {
	result := Database{log: log}

	err := result.NewConnection(ctx)
	if err != nil {
		log.Error().Msgf("[mongodb.New]Error. Can't connect to database. err = %s", err)
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

	DB.log.Info().Msg("[mongodb.NewConnection]Connected to MongoDB!")

	return nil
}

func (DB *Database) Disconnect(ctx context.Context) {
	if DB.c == nil {
		DB.log.Error().Msgf("[mongodb.Disconnect]Error. Database client not found")
		return
	}
	err := DB.c.Disconnect(ctx)
	if err != nil {
		DB.log.Error().Msgf("[mongodb.Disconnect]Error. Can't disconnect from database. err = %s", err)
		return
	}

	DB.log.Debug().Msg("[mongodb.Disconnect]Connection to MongoDB is also closed.")
}

func (DB *Database) NotFoundError() string {
	return ErrorNotExist
}
