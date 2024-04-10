package mongodb

import (
	"context"
	"fmt"
	"pidorator-bot/app/bot/trigger"
	"pidorator-bot/app/database"
	"pidorator-bot/app/database/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (DB *Database) GetBotData(ctx context.Context, guildID string) (*model.BotData, error) {
	result := &model.BotData{}

	if DB.botData == nil {
		DB.log.Debug().Msg("[mongodb.GetBotData]data map not found")
		err := DB.updateBotData(ctx)
		if err != nil {
			return result, err
		}
	}

	result, finded := DB.botData[guildID]
	if !finded {
		DB.log.Info().Msgf("[mongodb.GetBotData]data for guildID:%s not found. Creating new one", guildID)
		result, err := DB.createBotData(ctx, guildID)
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

	err = DB.updateBotData(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) createBotData(ctx context.Context, guildID string) (*model.BotData, error) {
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
	err = DB.updateBotData(ctx)
	if err != nil {
		return &model.BotData{}, err
	}

	return data, nil
}

func (DB *Database) updateBotData(ctx context.Context) error {
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
			DB.log.Error().Err(err).Msg("[mongodb.updateBotData.cur.next]")
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return err
	}
	DB.log.Debug().Msgf("[mongodb.updateBotData]Found multiple documents(%d). Array of pointers: %+v", len(results), results)

	// adds getted data into map
	for _, data := range results {
		DB.botData[data.GuildID] = data
	}

	return nil
}
