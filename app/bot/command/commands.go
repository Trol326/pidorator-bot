package command

import (
	"context"
	"pidorator-bot/app/database"
	"pidorator-bot/app/database/model"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

const (
	ErrorNoPermission string = "HTTP 403 Forbidden, {\"message\": \"Missing Permissions\", \"code\": 50013}"
	ErrorNoAccess     string = "HTTP 403 Forbidden, {\"message\": \"Missing Access\", \"code\": 50001}"
)

// Interface of commands for rolls
type Game interface {
	// main command for the game
	// gets timestamp of last "roll" from db
	// if currentTime - timestamp > 24h
	// do new "roll"
	// else
	// say to user that you can't do that
	Who(сtx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) (*model.EventData, error)
	// disable autoroll in game
	ChangeAutoRoll(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate)
	// adds new player in game db
	// TODO make it possible to add users by mention and UserID
	AddPlayer(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate)
	// writes full list of partisipians
	// TODO top 10, top 5
	List(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate)
	// writes full list of events for this server
	EventList(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate)
	// updates data for each player, such as local username
	UpdatePlayersData(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate)
	// TODO change season
	// StartNewSeason(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate)
}

// Interface of commands for administration purposes
type Admin interface {
	// Changes bot global name
	BotRename(сtx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate)

	// Changes bot prefix for this server
	SetPrefix(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate, newPrefix string)
}

// Contains bot command implementations
type Commands struct {
	log *zerolog.Logger
	db  database.Database
}

func New(l *zerolog.Logger, db database.Database) *Commands {
	return &Commands{
		log: l,
		db:  db,
	}
}

func (c *Commands) GetBotData(ctx context.Context, guildID string) (*model.BotData, error) {
	data, err := c.db.GetBotData(ctx, guildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't get bot data")
		return nil, err
	}

	return data, nil
}
