package trigger

import (
	"pidorator-bot/bot/command"
	"pidorator-bot/database"
	"pidorator-bot/utils/timer"

	"github.com/rs/zerolog"
)

type Trigger struct {
	Log      *zerolog.Logger
	commands *command.Commands
	timers   *timer.MasterTimer
}

func New(log *zerolog.Logger, db database.Database, masterTimer *timer.MasterTimer) *Trigger {
	c := command.New(log, db)

	return &Trigger{
		Log:      log,
		commands: c,
		timers:   masterTimer,
	}
}
