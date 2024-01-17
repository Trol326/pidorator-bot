package trigger

import (
	"context"
	"pidorator-bot/bot/command"
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
*/
func (t *Trigger) OnNewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID || !strings.HasPrefix(message.Content, DefaultPrefix) {
		return
	}

	t.Log.Debug().Msgf("Getted %s", message.Content)

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
	case strings.HasPrefix(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	case strings.HasPrefix(message.Content, "!bye"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
	case strings.HasPrefix(message.Content, "!ктопидор"):
		if game != nil {
			game.Who(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, "!botrename"):
		if admin != nil {
			admin.BotRename(ctx, discord, message)
		}
	default:
		t.Log.Debug().Msg("Command not found")
	}
}
