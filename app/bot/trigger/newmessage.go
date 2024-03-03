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
	if message.Author.ID == discord.State.User.ID || !strings.HasPrefix(message.Content, DefaultPrefix) {
		return
	}

	t.Log.Info().Msgf("Got %s", message.Content)

	ctx := context.Background()

	var game command.Game
	var admin command.Admin

	if GameEnabled {
		game = t.commands
	}
	if AdminEnabled {
		admin = t.commands
	}

	switch {
	case strings.HasPrefix(message.Content, DefaultPrefix+"help"):
		// TODO autodoc for help command
		discord.ChannelMessageSend(message.ChannelID, "Документация будет позже :D")
	case strings.HasPrefix(message.Content, DefaultPrefix+"ктопидор"):
		if game != nil {
			event, err := game.Who(ctx, discord, message)
			if err != nil {
				return
			}
			t.OnEventCreation(ctx, discord, event)
		}
	case strings.HasPrefix(message.Content, DefaultPrefix+"списокпидоров") || strings.HasPrefix(message.Content, DefaultPrefix+"топпидоров"):
		if game != nil {
			game.List(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, DefaultPrefix+"пидордня"):
		if game != nil {
			game.AddPlayer(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, DefaultPrefix+"обновитьпидоров"):
		if game != nil {
			game.UpdatePlayersData(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, DefaultPrefix+"списоксобытий"):
		if game != nil {
			game.EventList(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, DefaultPrefix+"botrename"):
		if admin != nil {
			admin.BotRename(ctx, discord, message)
		}
	default:
		t.Log.Debug().Msg("Command not found")
	}
}
