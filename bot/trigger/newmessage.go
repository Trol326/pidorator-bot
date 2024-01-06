package trigger

import (
	"pidorator-bot/bot/command"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const DefaultPrefix string = "!"

/*
Trigger when user send new message on server. Only in available for bot channels
*/
func (t *Trigger) OnNewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID || !strings.HasPrefix(message.Content, DefaultPrefix) {
		return
	}

	c := command.New(discord, message, t.Log)

	var game command.Game = &c
	var admin command.Admin = &c

	switch {
	case strings.HasPrefix(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	case strings.HasPrefix(message.Content, "!bye"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
	case strings.HasPrefix(message.Content, "!ктопидор"):
		game.Who()
	case strings.HasPrefix(message.Content, "!botrename"):
		admin.BotRename()
	}
}
