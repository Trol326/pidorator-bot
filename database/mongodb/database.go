package mongodb

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
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
		return Database{}
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

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	DB.log.Info().Msg("Connected to MongoDB!")

	return nil
}

func (DB *Database) Disconnect() {
	err := DB.c.Disconnect(context.TODO())

	if err != nil {
		DB.log.Error().Msgf("Error. Can't disconnect from database. err = %s", err)
		return
	}

	fmt.Println("Connection to MongoDB closed.")
}
