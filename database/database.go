package database

import (
	"context"
)

type Database interface {
	NewConnection(context.Context) error
	Disconnect(context.Context)
	IncreaseUserScore(ctx context.Context, guildID string, userID string) error
	GetAllPlayers(ctx context.Context, guildID string) ([]*GameData, error)
	UpdateUsersData(ctx context.Context, guildID string, data []*GameData) error
}

type GameData struct {
	GuildID  string `bson:"guildID,omitempty"`
	UserID   string `bson:"userID,omitempty"`
	Username string `bson:"username,omitempty"`
	Score    int32  `bson:"score,omitempty"`
}

/*
func (d *GameData) String() string {
	return fmt.Sprintf("{GuildID: %s, UserID: %s, Username: %s, Score: %d}\n", d.GuildID, d.UserID, d.Username, d.Score)
}
*/
