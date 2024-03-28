package trigger

import (
	"context"
	"pidorator-bot/app/bot/command"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	DefaultPrefix string = "!"
	GameEnabled   bool   = true
	AdminEnabled  bool   = true
)

/*
Trigger when user send new message on server. Only in available for bot channels

TODO maybe make trigger for game and for admin(?)
*/
func (t *Trigger) OnNewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}

	var prefix string
	ctx := context.Background()
	data, err := t.commands.GetBotData(ctx, message.GuildID)
	if err != nil {
		t.Log.Err(err).Msg("[trigger.OnNewMessage]error on getbotdata")
		prefix = DefaultPrefix
	} else {
		prefix = data.BotPrefix
	}

	if !strings.HasPrefix(message.Content, prefix) {
		return
	}

	t.Log.Info().Msgf("Got %s", message.Content)

	var game command.Game
	var admin command.Admin

	if GameEnabled {
		game = t.commands
	}
	if AdminEnabled {
		admin = t.commands
	}

	switch {
	case strings.HasPrefix(message.Content, prefix+"help"):
		// TODO autodoc for help command
		_, err := discord.ChannelMessageSend(message.ChannelID, "Документация будет позже :D")
		if err != nil {
			t.Log.Err(err).Msg("[trigger.OnNewMessage]error on channelMessageSend")
			return
		}
	case strings.HasPrefix(message.Content, prefix+"ктопидор"):
		if game != nil {
			event, err := game.Who(ctx, discord, message)
			if err != nil {
				return
			}
			t.OnEventCreation(ctx, discord, event)
		}
	case strings.HasPrefix(message.Content, prefix+"списокпидоров") || strings.HasPrefix(message.Content, prefix+"топпидоров"):
		if game != nil {
			game.List(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, prefix+"пидордня"):
		if game != nil {
			game.AddPlayer(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, prefix+"обновитьпидоров"):
		if game != nil {
			game.UpdatePlayersData(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, prefix+"списоксобытий"):
		if game != nil {
			game.EventList(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, prefix+"botrename"):
		if admin != nil {
			admin.BotRename(ctx, discord, message)
		}
	default:
		t.Log.Debug().Msg("Command not found")
	}
}
