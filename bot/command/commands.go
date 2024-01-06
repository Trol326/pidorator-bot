package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

// Interface of commands for rolls
type Game interface {
	// you
	Who()
}

// Interface of commands for administration purposes
type Admin interface {
	// Changes bot global name
	BotRename()
}

// Contains bot command implementations
type Commands struct {
	discord *discordgo.Session
	message *discordgo.MessageCreate
	log     *zerolog.Logger
}

func New(d *discordgo.Session, m *discordgo.MessageCreate, l *zerolog.Logger) Commands {
	return Commands{
		discord: d,
		message: m,
		log:     l,
	}
}
