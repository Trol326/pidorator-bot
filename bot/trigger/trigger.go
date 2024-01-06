package trigger

import "github.com/rs/zerolog"

type Trigger struct {
	Log *zerolog.Logger
}

func New(log *zerolog.Logger) *Trigger {
	return &Trigger{
		Log: log,
	}
}
