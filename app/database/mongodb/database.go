package mongodb

import (
	"context"
	"fmt"
	"os"
	"pidorator-bot/app/bot/trigger"
	"pidorator-bot/app/database"
	"pidorator-bot/app/database/model"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
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

// TODO refactor functions to different files

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

func (DB *Database) UpdateBotData(ctx context.Context) error {
	if DB.botData == nil {
		DB.botData = make(map[string]*model.BotData)
	}

	// Get all bot data from DB
	c := DB.c.Database(BotDBName).Collection(DataCollectionName)
	findOptions := options.Find()
	var results []*model.BotData
	cur, err := c.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem model.BotData
		err := cur.Decode(&elem)
		if err != nil {
			DB.log.Error().Err(err).Msg("[mongodb.UpdateBotData.cur.next]")
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return err
	}
	DB.log.Debug().Msgf("[mongodb.UpdateBotData]Found multiple documents(%d). Array of pointers: %+v", len(results), results)

	// adds getted data into map
	for _, data := range results {
		DB.botData[data.GuildID] = data
	}

	return nil
}

func (DB *Database) GetBotData(ctx context.Context, guildID string) (*model.BotData, error) {
	result := &model.BotData{}

	if DB.botData == nil {
		DB.log.Debug().Msg("[mongodb.GetBotData]data map not found")
		err := DB.UpdateBotData(ctx)
		if err != nil {
			return result, err
		}
	}

	result, finded := DB.botData[guildID]
	if !finded {
		DB.log.Info().Msgf("[mongodb.GetBotData]data for guildID:%s not found. Creating new one", guildID)
		result, err := DB.CreateBotData(ctx, guildID)
		if err != nil {
			return result, err
		}
		return result, nil
	}

	return result, nil
}

func (DB *Database) ChangeBotData(ctx context.Context, data *model.BotData) error {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return err
	}

	opts := options.Replace().SetUpsert(true)
	c := DB.c.Database(BotDBName).Collection(DataCollectionName)
	filter := bson.D{{Key: "guildID", Value: data.GuildID}}

	result, err := c.ReplaceOne(ctx, filter, data, opts)
	if err != nil {
		return err
	}

	DB.log.Debug().Msgf("[mongodb.ChangeBotData]Matched %v documents and updated %v documents.", result.MatchedCount, result.ModifiedCount)

	err = DB.UpdateBotData(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) CreateBotData(ctx context.Context, guildID string) (*model.BotData, error) {
	// Default bot data
	data := &model.BotData{
		GuildID:           guildID,
		IsGameEnabled:     true,
		IsAutoRollEnabled: true,
		BotPrefix:         trigger.DefaultPrefix,
	}

	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return &model.BotData{}, err
	}

	c := DB.c.Database(BotDBName).Collection(DataCollectionName)

	findOptions := options.FindOne()
	filter := bson.D{{Key: "guildID", Value: data.GuildID}}

	d := model.BotData{}
	err := c.FindOne(ctx, filter, findOptions).Decode(&d)
	if err != nil && err.Error() != ErrorNotExist {
		return &model.BotData{}, err
	}
	if d.GuildID != "" {
		err := fmt.Errorf(database.DataAlreadyExistError)
		return &model.BotData{}, err
	}

	result, err := c.InsertOne(ctx, data)
	if err != nil {
		return &model.BotData{}, err
	}

	if result == nil {
		err := fmt.Errorf("error. Data not created")
		return &model.BotData{}, err
	}

	// Update map
	err = DB.UpdateBotData(ctx)
	if err != nil {
		return &model.BotData{}, err
	}

	return data, nil
}

func (DB *Database) AddPlayer(ctx context.Context, data *model.PlayerData) error {
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

	result, err := c.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	if result == nil {
		err := fmt.Errorf("error. Player not added")
		return err
	}

	return nil
}

func (DB *Database) AddEvent(ctx context.Context, data *model.EventData) error {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return err
	}

	p, err := DB.FindEvent(ctx, data.GuildID, data.Type)
	if err != nil && err.Error() != ErrorNotExist {
		return err
	}

	if p.GuildID != "" {
		err := fmt.Errorf(database.EventAlreadyExistError)
		return err
	}

	c := DB.c.Database(BotDBName).Collection(EventCollectionName)

	result, err := c.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	if result == nil {
		err := fmt.Errorf("error. Event not added")
		return err
	}

	DB.log.Debug().Msgf("[mongodb.AddEvent]Event added successfully: %v", result)

	return nil
}

func (DB *Database) UpdateEvent(ctx context.Context, data *model.EventData) error {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return err
	}

	c := DB.c.Database(BotDBName).Collection(EventCollectionName)
	filter := bson.D{{Key: "guildID", Value: data.GuildID}, {Key: "eventType", Value: data.Type}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "startTime", Value: data.StartTime}, {Key: "endTime", Value: data.EndTime}, {Key: "channelID", Value: data.ChannelID}}}}

	result, err := c.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	DB.log.Debug().Msgf("[mongodb.UpdateEvent]Matched %v documents and updated %v documents.", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (DB *Database) GetEvents(ctx context.Context, guildID string) ([]*model.EventData, error) {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return []*model.EventData{}, err
	}

	c := DB.c.Database(BotDBName).Collection(EventCollectionName)

	findOptions := options.Find().SetSort(bson.D{{Key: "startTime", Value: -1}})

	var results []*model.EventData

	filter := bson.D{}
	if guildID != "" {
		filter = bson.D{{Key: "guildID", Value: guildID}}
	}

	cur, err := c.Find(ctx, filter, findOptions)
	if err != nil {
		return []*model.EventData{}, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem model.EventData
		err := cur.Decode(&elem)
		if err != nil {
			DB.log.Error().Err(err).Msg("[mongodb.GetEvents.cur.next]")
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return []*model.EventData{}, err
	}

	DB.log.Debug().Msgf("[mongodb.GetEvents]Found multiple documents(%d). Array of pointers: %+v", len(results), results)

	return results, nil
}

func (DB *Database) GetAllPlayers(ctx context.Context, guildID string, sortingType int) ([]*model.PlayerData, error) {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return []*model.PlayerData{}, err
	}

	c := DB.c.Database(BotDBName).Collection(GameCollectionName)

	findOptions := options.Find()
	if sortingType != database.NoSorting {
		findOptions = findOptions.SetSort(bson.D{{Key: "score", Value: sortingType}})
	}

	var results []*model.PlayerData

	cur, err := c.Find(ctx, bson.D{{Key: "guildID", Value: guildID}}, findOptions)
	if err != nil {
		return []*model.PlayerData{}, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem model.PlayerData
		err := cur.Decode(&elem)
		if err != nil {
			DB.log.Error().Err(err).Msg("[mongodb.GetAllPlayers.cur.next]")
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return []*model.PlayerData{}, err
	}

	DB.log.Debug().Msgf("[mongodb.GetAllPlayers]Found multiple documents(%d). Array of pointers: %+v", len(results), results)

	return results, nil
}

func (DB *Database) FindPlayer(ctx context.Context, guildID string, userID string) (*model.PlayerData, error) {
	result := model.PlayerData{}
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return &result, err
	}

	c := DB.c.Database(BotDBName).Collection(GameCollectionName)

	findOptions := options.FindOne()

	filter := bson.D{{Key: "guildID", Value: guildID}, {Key: "userID", Value: userID}}

	err := c.FindOne(ctx, filter, findOptions).Decode(&result)
	if err != nil {
		return &model.PlayerData{}, err
	}
	DB.log.Debug().Msgf("[mongodb.FindPlayer]Player finded: %v", result)

	return &result, nil
}

func (DB *Database) FindEvent(ctx context.Context, guildID string, eventType string) (*model.EventData, error) {
	result := model.EventData{}
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return &result, err
	}

	c := DB.c.Database(BotDBName).Collection(EventCollectionName)

	findOptions := options.FindOne()

	filter := bson.D{{Key: "guildID", Value: guildID}, {Key: "eventType", Value: eventType}}

	err := c.FindOne(ctx, filter, findOptions).Decode(&result)
	if err != nil {
		return &model.EventData{}, err
	}

	DB.log.Debug().Msgf("[mongodb.FindEvent]Event finded: %v", result)

	return &result, nil
}

func (DB *Database) IncreasePlayerScore(ctx context.Context, guildID string, userID string) error {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return err
	}

	c := DB.c.Database(BotDBName).Collection(GameCollectionName)
	filter := bson.D{{Key: "userID", Value: userID}, {Key: "guildID", Value: guildID}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "score", Value: 1}}}}

	result, err := c.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	DB.log.Debug().Msgf("[mongodb.IncreasePlayerScore]Matched %v documents and updated %v documents.", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (DB *Database) UpdatePlayersData(ctx context.Context, data []*model.PlayerData) error {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return err
	}

	DB.log.Debug().Msgf("[mongodb.UpdatePlayersData]Got %d players data: %v", len(data), data)

	var counterMatched int64 = 0
	var counterModified int64 = 0
	c := DB.c.Database(BotDBName).Collection(GameCollectionName)
	for _, player := range data {
		filter := bson.D{{Key: "userID", Value: player.UserID}, {Key: "guildID", Value: player.GuildID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "username", Value: player.Username}}}}

		result, err := c.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		counterMatched += result.MatchedCount
		counterModified += result.ModifiedCount
	}

	DB.log.Debug().Msgf("[mongodb.UpdatePlayersData]Matched %v documents and updated %v documents", counterMatched, counterModified)

	return nil
}

func (DB *Database) NotFoundError() string {
	return ErrorNotExist
}
