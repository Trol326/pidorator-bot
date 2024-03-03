package command

import (
	"context"
	"pidorator-bot/app/database"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

// Interface of commands for rolls
type Game interface {
	// main command for the game
	// gets timestamp of last "roll" from db
	// if currentTime - timestamp > 24h
	// do new "roll"
	// else
	// say to user that you can't do that
	Who(сtx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) (*database.EventData, error)
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
