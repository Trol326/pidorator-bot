package mongodb

import (
	"context"
	"fmt"
	"pidorator-bot/app/database"
	"pidorator-bot/app/database/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
