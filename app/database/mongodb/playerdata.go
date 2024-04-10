package mongodb

import (
	"context"
	"fmt"
	"pidorator-bot/app/database"
	"pidorator-bot/app/database/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (DB *Database) AddPlayer(ctx context.Context, data *model.PlayerData) error {
	if DB.c == nil {
		err := fmt.Errorf("error. Database client not found")
		return err
	}

	p, err := DB.findPlayer(ctx, data.GuildID, data.UserID)
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

func (DB *Database) findPlayer(ctx context.Context, guildID string, userID string) (*model.PlayerData, error) {
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
	DB.log.Debug().Msgf("[mongodb.findPlayer]Player finded: %v", result)

	return &result, nil
}
