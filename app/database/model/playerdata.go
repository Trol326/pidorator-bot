package model

import "fmt"

type PlayerData struct {
	GuildID  string `bson:"guildID,omitempty"`
	UserID   string `bson:"userID,omitempty"`
	Username string `bson:"username,omitempty"`
	Score    int32  `bson:"score"`
}

func (d *PlayerData) String() string {
	return fmt.Sprintf("{GuildID: %s, UserID: %s, Username: %s, Score: %d}", d.GuildID, d.UserID, d.Username, d.Score)
}
