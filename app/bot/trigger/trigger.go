package trigger

import (
	"fmt"
	"pidorator-bot/app/bot/command"
	"pidorator-bot/app/database"
	"pidorator-bot/tools/timer"
	"regexp"

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

func ParseCommand(prefix, command string) ([]string, error) {
	p := regexp.QuoteMeta(prefix)
	exp := fmt.Sprintf("(%s[a-z]* )(.*)", p)
	regex, err := regexp.Compile(exp)
	if err != nil {
		return nil, err
	}
	return regex.FindStringSubmatch(command), nil
}
