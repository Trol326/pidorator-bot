package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (c *Commands) Who(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Debug().Msg("commands.who triggered")

	err := c.db.IncreaseUserScore(ctx, message.GuildID, message.Author.ID)
	if err != nil {
		c.log.Error().Err(err).Msgf("Error. Can't update, err = %s", err)
	}

	data, err := c.db.GetAllPlayers(ctx)
	if err != nil {
		c.log.Error().Err(err).Msgf("Error. Can't get players, err = %s", err)
	}

	c.log.Debug().Msgf("Getted %d players. Data: %v", len(data), data[0])

	text := fmt.Sprintf("%s, ты пидор <:MumeiYou:1192139708222935050>", message.Author.Mention())
	discord.ChannelMessageSend(message.ChannelID, text)
}
