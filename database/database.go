package database

import (
	"context"
	"fmt"
)

const (
	GameEventName string = "pidorator-timer"

	PlayerAlreadyExistError string = "error player already exist"
	EventAlreadyExistError  string = "error event already exist"
)

type Database interface {
	NewConnection(context.Context) error
	Disconnect(context.Context)
	AddPlayer(ctx context.Context, data *PlayerData) error
	GetAllPlayers(ctx context.Context, guildID string) ([]*PlayerData, error)
	// Adds new event, updates it if event already exist
	AddEvent(ctx context.Context, data *EventData) error
	UpdateEvent(ctx context.Context, data *EventData) error
	GetEvents(ctx context.Context, guildID string) ([]*EventData, error)
	FindEvent(ctx context.Context, guildID string, eventType string) (*EventData, error)
	IncreasePlayerScore(ctx context.Context, guildID string, userID string) error
	UpdatePlayersData(ctx context.Context, data []*PlayerData) error
}

//TODO models folder

type PlayerData struct {
	GuildID  string `bson:"guildID,omitempty"`
	UserID   string `bson:"userID,omitempty"`
	Username string `bson:"username,omitempty"`
	Score    int32  `bson:"score"`
}

type EventData struct {
	GuildID   string `bson:"guildID,omitempty"`
	Type      string `bson:"eventType,omitempty"`
	StartTime int64  `bson:"startTime"`
	EndTime   int64  `bson:"endTime"`
}

func (d *PlayerData) String() string {
	return fmt.Sprintf("{GuildID: %s, UserID: %s, Username: %s, Score: %d}", d.GuildID, d.UserID, d.Username, d.Score)
}

func (d *EventData) String() string {
	return fmt.Sprintf("{GuildID: %s, Type: %s, StartTime: %d, EndTime: %d}", d.GuildID, d.Type, d.StartTime, d.EndTime)
}

func (d EventData) IsEventEnded(now ...int64) bool {
	if d.StartTime > d.EndTime {
		return true
	}
	if d.StartTime == 0 || d.EndTime == 0 {
		return false
	}
	if len(now) < 1 {
		return false
	}
	return d.EndTime < now[0]
}
