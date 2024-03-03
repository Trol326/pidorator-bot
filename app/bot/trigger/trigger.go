package trigger

import (
	"pidorator-bot/app/bot/command"
	"pidorator-bot/app/database"
	"pidorator-bot/tools/timer"

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
