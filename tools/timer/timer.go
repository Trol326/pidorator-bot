package timer

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type MasterTimer struct {
	Timers map[string](*time.Timer)
	log    *zerolog.Logger
}

func New(l *zerolog.Logger) MasterTimer {
	result := MasterTimer{
		Timers: make(map[string]*time.Timer, 0),
		log:    l,
	}
	return result
}

// Returns timer object and final timername
func (t *MasterTimer) New(name string, secs int64) (timer *time.Timer, timerName string) {
	if secs < 0 {
		return nil, ""
	}

	timer = time.NewTimer(time.Second * time.Duration(secs))

	timerName = name
	if _, ok := t.Timers[timerName]; ok {
		t.log.Debug().Msgf("%s already exist", name)
		for i := 0; ; i++ {
			timerName = fmt.Sprintf("%s_%d", name, i)
			_, ok := t.Timers[timerName]
			if !ok {
				break
			}
			t.log.Debug().Msgf("%s still already exist", timerName)
			continue
		}
	}

	t.Timers[timerName] = timer

	return timer, timerName
}

func (t *MasterTimer) StopAll() {
	for k, v := range t.Timers {
		go func() {
			if !v.Stop() {
				t.log.Info().Msgf("Timer \"%s\" ended", k)
				<-v.C
			}
		}()
		t.log.Info().Msgf("Timer \"%s\" was stopped", k)
	}
	t.Timers = make(map[string]*time.Timer)
	t.log.Info().Msgf("All timers were stopped")
}

func (t *MasterTimer) StopByName(name string) {
	timer, ok := t.Timers[name]
	if !ok {
		t.log.Info().Msgf("Timer \"%s\" not found", name)
	}

	go func() {
		if !timer.Stop() {
			t.log.Debug().Msgf("Timer \"%s\" ended", name)
			<-timer.C
		} else {
			t.log.Debug().Msgf("Timer \"%s\" stopped", name)
		}
	}()
	delete(t.Timers, name)
	t.log.Info().Msgf("Timer \"%s\" deleted", name)
}
