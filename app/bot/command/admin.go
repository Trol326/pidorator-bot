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

func (c *Commands) SetPrefix(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate, newPrefix string) {
	c.log.Info().Msg("[commands.SetPrefix]triggered")

	data, err := c.db.GetBotData(ctx, message.GuildID)
	if err != nil {
		_, err = discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if err != nil {
			c.log.Err(err).Msg("[commands.SetPrefix]error on channelMessageSend")
		}
	}

	data.BotPrefix = newPrefix

	err = c.db.ChangeBotData(ctx, data)
	if err != nil {
		_, err = discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if err != nil {
			c.log.Err(err).Msg("[commands.SetPrefix]error on channelMessageSend")
		}
		return
	}

	text := fmt.Sprintf("Префикс успешно изменен на, `%s`", newPrefix)
	_, err = discord.ChannelMessageSend(message.ChannelID, text)
	if err != nil {
		c.log.Err(err).Msg("[commands.SetPrefix]error on channelMessageSend")
	}
}
