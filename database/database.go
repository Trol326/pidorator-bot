package database

import (
	"context"
)

const (
	PlayerAlreadyExistError string = "error player already exist"
)

type Database interface {
	NewConnection(context.Context) error
	Disconnect(context.Context)
	AddPlayer(ctx context.Context, data *PlayerData) error
	GetAllPlayers(ctx context.Context, guildID string) ([]*PlayerData, error)
	IncreasePlayerScore(ctx context.Context, guildID string, userID string) error
	UpdatePlayersData(ctx context.Context, guildID string, data []*PlayerData) error
}

type PlayerData struct {
	GuildID  string `bson:"guildID,omitempty"`
	UserID   string `bson:"userID,omitempty"`
	Username string `bson:"username,omitempty"`
	Score    int32  `bson:"score"`
}

/*
func (d *GameData) String() string {
	return fmt.Sprintf("{GuildID: %s, UserID: %s, Username: %s, Score: %d}\n", d.GuildID, d.UserID, d.Username, d.Score)
}
*/
