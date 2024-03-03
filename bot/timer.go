package bot

import (
	"context"
	"fmt"
	"pidorator-bot/database"
)

func (c *Client) InitTimers(ctx context.Context) {
	if c.DB == nil {
		return
	}

	c.Log.Info().Msgf("Starting event timers...")
	events, err := c.DB.GetEvents(ctx, "")
	if err != nil {
		c.Log.Error().Err(err).Msg("[bot.InitTimers]")
		return
	}

	for _, event := range events {
		if event.Type != database.GameEventName {
			continue
		}

		text := fmt.Sprintf("timer_%s_%s", event.Type, event.GuildID)
		t, name := c.Timers.New(text, event.SecondsUntilEnd())
		go func() {
			<-t.C
			c.Triggers.OnTimerEnded(ctx, c.Session, event.GuildID, event.ChannelID, event.Type)
			c.Timers.StopByName(name)
		}()
		c.Log.Info().Msgf("Started %s", text)
	}
}
