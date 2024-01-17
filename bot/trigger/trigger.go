package trigger

import (
	"pidorator-bot/bot/command"
	"pidorator-bot/database"

	"github.com/rs/zerolog"
)

type Trigger struct {
	Log      *zerolog.Logger
	commands *command.Commands
}

func New(log *zerolog.Logger, db database.Database) *Trigger {
	c := command.New(log, db)

	return &Trigger{
		Log:      log,
		commands: c,
	}
}
