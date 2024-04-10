package database

import (
	"context"
	"pidorator-bot/app/database/model"
)

const (
	GameEventName string = "pidorator-timer"

	PlayerAlreadyExistError string = "error player already exist"
	EventAlreadyExistError  string = "error event already exist"
	DataAlreadyExistError   string = "error data already exist"

	AccendingSorting int = 1
	DcendingSorting  int = -1
	NoSorting        int = 0
)

type Database interface {

	/* 	Connection */

	// Creates new connection to db and ping it
	NewConnection(context.Context) error
	// Destroys connection. Use in defer near to NewConnection
	Disconnect(context.Context)

	/* Bot Data */

	// Returns config data for this guildID
	GetBotData(ctx context.Context, guildID string) (*model.BotData, error)
	// Updates bot data in db, creates new one if it doesn't exist
	ChangeBotData(ctx context.Context, data *model.BotData) error

	/* Events */

	// Creates new event, updates it if event already exist
	AddEvent(ctx context.Context, data *model.EventData) error
	// Updates event in database
	UpdateEvent(ctx context.Context, data *model.EventData) error
	// Returns unsorted array of events for this server (guildID)
	GetEvents(ctx context.Context, guildID string) ([]*model.EventData, error)
	// Finds specific event based on eventType and guildID
	FindEvent(ctx context.Context, guildID string, eventType string) (*model.EventData, error)

	/* Pidorator Game */

	// Creates new player in db
	AddPlayer(ctx context.Context, data *model.PlayerData) error
	// Returns array of players, sorting based on sortingType. Consts AccendingSorting/DcendingSorting/NoSorting
	GetAllPlayers(ctx context.Context, guildID string, sortingType int) ([]*model.PlayerData, error)
	// Increases player score at one
	IncreasePlayerScore(ctx context.Context, guildID string, userID string) error
	// Updates all players username in db
	UpdatePlayersData(ctx context.Context, data []*model.PlayerData) error

	/* Utility */

	// Returns driver specific error string
	NotFoundError() string
}
