package database

import "context"

type Database interface {
	NewConnection(context.Context) error
	Disconnect(context.Context)
	IncreaseUserScore(ctx context.Context, guildID string, userID string) error
	GetAllPlayers(ctx context.Context) ([]*GameData, error)
}

type GameData struct {
	GuildID string `bson:"guildID,omitempty"`
	UserID  string `bson:"userID,omitempty"`
	Score   int32  `bson:"score,omitempty"`
}
