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
		c.log.Error().Err(err).Msgf("Error. Can't update")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	text := fmt.Sprintf("%s, ты пидор <:MumeiYou:1192139708222935050>", message.Author.Mention())
	discord.ChannelMessageSend(message.ChannelID, text)
}

func (c *Commands) AddPlayer(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {

}

func (c *Commands) List(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Debug().Msg("commands.list triggered")

	data, err := c.db.GetAllPlayers(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("Error. Can't get players")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	c.log.Debug().Msgf("Getted %d players", len(data))

	top := ""
	for i, d := range data {
		if d.Username == "" {
			top += fmt.Sprintf("%d) %s - %d\n", i+1, d.UserID, d.Score)
			continue
		}
		top += fmt.Sprintf("%d) %s - %d\n", i+1, d.Username, d.Score)
	}

	text := fmt.Sprintf("Топ пидоров:\n%s", top)
	discord.ChannelMessageSend(message.ChannelID, text)
}
