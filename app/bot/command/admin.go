package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// TODO make some functionality
func (c *Commands) BotRename(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Debug().Msg("commands.botrename triggered")
	text := fmt.Sprintf("Ты пидор, %s <:MumeiYou:1192139708222935050>", message.Author.Mention())
	_, err := discord.ChannelMessageSend(message.ChannelID, text)
	if err != nil {
		c.log.Err(err).Msg("[commands.botrename]error on channelMessageSend")
	}
}
