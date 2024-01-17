package command

import (
	"context"
	"pidorator-bot/database"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

// Interface of commands for rolls
type Game interface {
	// you
	Who(сtx context.Context, d *discordgo.Session, m *discordgo.MessageCreate)
}

// Interface of commands for administration purposes
type Admin interface {
	// Changes bot global name
	BotRename(сtx context.Context, d *discordgo.Session, m *discordgo.MessageCreate)
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
