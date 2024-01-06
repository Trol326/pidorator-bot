package command

import (
	"fmt"
)

func (c *Commands) BotRename() {
	c.log.Debug().Msg("commands.botrename triggered")
	text := fmt.Sprintf("Ты пидор, %s <:MumeiYou:1192139708222935050>", c.message.Author.Mention())
	c.discord.ChannelMessageSend(c.message.ChannelID, text)
}
