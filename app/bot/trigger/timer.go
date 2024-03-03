package trigger

import (
	"context"
	"fmt"
	"pidorator-bot/app/database"

	"github.com/bwmarrin/discordgo"
)

func (t *Trigger) OnTimerEnded(ctx context.Context, discord *discordgo.Session, guildID string, channelID string, timerType string) {
	if channelID == "" {
		t.Log.Error().Msgf("[Trigger.OnTimerEnded]Timer ChannelID not found. GuildID = %s; TimerType = %s", guildID, timerType)
		return
	}

	if timerType == database.GameEventName && GameEnabled {
		text := "Время пришло. Пора выбирать пидора!"
		discord.ChannelMessageSend(channelID, text)
		event, err := t.commands.AutoRoll(ctx, discord, guildID, channelID)
		if err != nil {
			return
		}
		t.OnEventCreation(ctx, discord, event)
	}
}

func (t *Trigger) OnEventCreation(ctx context.Context, discord *discordgo.Session, event *database.EventData) {
	t.Log.Info().Msgf("Starting event timer...")
	text := fmt.Sprintf("timer_%s_%s", event.Type, event.GuildID)
	timer, name := t.timers.New(text, event.SecondsUntilEnd())
	go func() {
		<-timer.C
		t.OnTimerEnded(ctx, discord, event.GuildID, event.ChannelID, event.Type)
		t.timers.StopByName(name)
	}()
	t.Log.Info().Msgf("Started %s", name)
}
