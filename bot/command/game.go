package command

import (
	"context"
	"fmt"
	"math/rand"
	"pidorator-bot/database"
	"pidorator-bot/utils"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (c *Commands) Who(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Debug().Msg("commands.who triggered")

	ev, err := c.db.FindEvent(ctx, message.GuildID, database.GameEventName)
	if err != nil && err.Error() != database.EventAlreadyExistError {
		c.log.Error().Err(err).Msgf("Error. Can't find event")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}
	if ev.GuildID != "" && !IsGameEventEnded(ev) {
		c.log.Debug().Msg("Can't roll, event still not ended")
		text := fmt.Sprintf("Да ты погоди, %s, время ещё не пришло. В предыдущий раз крутили %s", message.Author.GlobalName, utils.ToDiscordTimeStamp(ev.StartTime, utils.TSFormat().Relative))
		discord.ChannelMessageSend(message.ChannelID, text)
		return
	}

	player, err := c.getRandomPlayer(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("Error. Can't get player")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	err = c.db.IncreasePlayerScore(ctx, player.GuildID, player.UserID)
	if err != nil {
		c.log.Error().Err(err).Msgf("Error. Can't increase player score")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Hour * 24)
	event := database.EventData{
		GuildID:   message.GuildID,
		Type:      database.GameEventName,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
	}
	err = c.db.AddEvent(ctx, &event)
	if err != nil && err.Error() != database.EventAlreadyExistError {
		c.log.Error().Err(err).Msgf("Error. Can't add event")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}
	if err.Error() == database.EventAlreadyExistError {
		err := c.db.UpdateEvent(ctx, &event)
		if err != nil {
			c.log.Error().Err(err).Msgf("Error. Can't update event")
			discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		}
	}

	text := fmt.Sprintf("%s, ты пидор <:MumeiYou:1192139708222935050>", utils.UserIDToMention(player.UserID))
	discord.ChannelMessageSend(message.ChannelID, text)
}

func (c *Commands) AddPlayer(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Debug().Msg("commands.addplayer triggered")

	nickname := message.Member.Nick
	if nickname == "" {
		// TODO for discordgo pre 0.27.2, delete when package updates
		//nickname = cases.Title(language.Und, cases.NoLower).String(message.Author.Username)
		nickname = message.Author.GlobalName
	}
	data := database.PlayerData{GuildID: message.GuildID, UserID: message.Author.ID, Username: nickname}

	err := c.db.AddPlayer(ctx, &data)
	if err != nil {
		if err.Error() == database.PlayerAlreadyExistError {
			c.log.Debug().Err(err).Msgf("player already exist")
			text := fmt.Sprintf("э, %s, куда тебе? Ты же уже участвуешь", data.Username)
			discord.ChannelMessageSend(message.ChannelID, text)
			return
		}
		c.log.Error().Err(err).Msgf("Error. Can't add player")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	text := fmt.Sprintf("%s, вы приняты в пидорскую рулетку", data.Username)
	discord.ChannelMessageSend(message.ChannelID, text)
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

func (c *Commands) EventList(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Debug().Msg("commands.eventList triggered")

	data, err := c.db.GetEvents(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("Error. Can't get events")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	c.log.Debug().Msgf("Getted %d events", len(data))

	top := ""
	for i, d := range data {
		startTS := utils.ToDiscordTimeStamp(d.StartTime, utils.TSFormat().LongDateShortTime)
		endTS := utils.ToDiscordTimeStamp(d.EndTime, utils.TSFormat().LongDateShortTime)
		top += fmt.Sprintf("%d) %s. %s - %s\n", i+1, d.Type, startTS, endTS)
	}

	text := fmt.Sprintf("Список запланированных на этом сервере событий:\n%s", top)
	discord.ChannelMessageSend(message.ChannelID, text)
}

// Checks is game event ended by time. True if ended, False otherwise
func IsGameEventEnded(event *database.EventData) bool {
	now := time.Now()
	res := event.IsEventEnded(now.Unix())
	if res {
		return res
	}

	endTime := time.Unix(event.StartTime, 0).Add(time.Hour * 24)

	// if startTime + 24h < now => returns true
	return endTime.Before(now)
}

func (c *Commands) getRandomPlayer(ctx context.Context, guildID string) (*database.PlayerData, error) {
	result := &database.PlayerData{}

	players, err := c.db.GetAllPlayers(ctx, guildID)
	if err != nil {
		return &database.PlayerData{}, err
	}

	i := len(players) - 1
	random := rand.New(rand.NewSource(time.Now().Unix()))
	num := random.Int31n(int32(i))

	result = players[num]

	return result, nil
}
