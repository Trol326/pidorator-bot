package command

import (
	"fmt"
)

func (c *Commands) Who() {
	c.log.Debug().Msg("commands.who triggered")
	text := fmt.Sprintf("%s, ты пидор <:MumeiYou:1192139708222935050>", c.message.Author.Mention())
	c.discord.ChannelMessageSend(c.message.ChannelID, text)
}
