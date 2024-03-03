package command

import (
	"context"
	"fmt"
	"math/rand"
	"pidorator-bot/app/database"
	"pidorator-bot/tools"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (c *Commands) AutoRoll(ctx context.Context, discord *discordgo.Session, guildID string, channelID string) (*database.EventData, error) {
	c.log.Info().Msg("[commands.Autoroll]triggered")

	ev, err := c.db.FindEvent(ctx, guildID, database.GameEventName)
	if err != nil && err.Error() != database.EventAlreadyExistError {
		c.log.Error().Err(err).Msgf("[commands.Autoroll]Error. Can't find event")
		discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
		return nil, err
	}
	if ev.GuildID != "" && !IsGameEventEnded(ev) {
		c.log.Debug().Msg("[commands.Autoroll]Can't roll, event still not ended")
		text := fmt.Sprintf("Накосячил, ещё рано. В предыдущий раз крутили %s", tools.ToDiscordTimeStamp(ev.StartTime, tools.TSFormat().Relative))
		discord.ChannelMessageSend(channelID, text)
		err := fmt.Errorf("event not ended")
		return nil, err
	}

	player, err := c.getRandomPlayer(ctx, guildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.Autoroll]Error. Can't get player")
		discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
		return nil, err
	}

	err = c.db.IncreasePlayerScore(ctx, player.GuildID, player.UserID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.Autoroll]Error. Can't increase player score")
		discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
		return nil, err
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Hour * 24)
	event := database.EventData{
		GuildID:   guildID,
		ChannelID: channelID,
		Type:      database.GameEventName,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
	}
	err = c.db.AddEvent(ctx, &event)
	if err != nil && err.Error() != database.EventAlreadyExistError {
		c.log.Error().Err(err).Msgf("[commands.Autoroll]Error. Can't add event")
		discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
		return nil, err
	}
	if err.Error() == database.EventAlreadyExistError {
		err := c.db.UpdateEvent(ctx, &event)
		if err != nil {
			c.log.Error().Err(err).Msgf("[commands.Autoroll]Error. Can't update event")
			discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
			return nil, err
		}
	}

	text := fmt.Sprintf("%s, ты пидор <:MumeiYou:1192139708222935050>", tools.UserIDToMention(player.UserID))
	discord.ChannelMessageSend(channelID, text)
	return &event, nil
}

func (c *Commands) Who(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) (*database.EventData, error) {
	c.log.Info().Msg("[commands.Who]triggered")

	ev, err := c.db.FindEvent(ctx, message.GuildID, database.GameEventName)
	if err != nil && err.Error() != database.EventAlreadyExistError {
		c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't find event")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return nil, err
	}
	if ev.GuildID != "" && !IsGameEventEnded(ev) {
		c.log.Debug().Msg("[commands.Who]Can't roll, event still not ended")
		text := fmt.Sprintf("Да ты погоди, %s, время ещё не пришло. В предыдущий раз крутили %s", message.Author.GlobalName, tools.ToDiscordTimeStamp(ev.StartTime, tools.TSFormat().Relative))
		discord.ChannelMessageSend(message.ChannelID, text)
		err = fmt.Errorf("event not ended")
		return nil, err
	}

	player, err := c.getRandomPlayer(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't get player")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return nil, err
	}

	err = c.db.IncreasePlayerScore(ctx, player.GuildID, player.UserID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't increase player score")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return nil, err
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Hour * 24)
	event := database.EventData{
		GuildID:   message.GuildID,
		ChannelID: message.ChannelID,
		Type:      database.GameEventName,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
	}
	err = c.db.AddEvent(ctx, &event)
	if err != nil && err.Error() != database.EventAlreadyExistError {
		c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't add event")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return nil, err
	}
	if err.Error() == database.EventAlreadyExistError {
		err := c.db.UpdateEvent(ctx, &event)
		if err != nil {
			c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't update event")
			discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
			return nil, err
		}
	}

	text := fmt.Sprintf("%s, ты пидор <:MumeiYou:1192139708222935050>", tools.UserIDToMention(player.UserID))
	discord.ChannelMessageSend(message.ChannelID, text)
	return &event, nil
}

func (c *Commands) AddPlayer(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Info().Msg("[commands.AddPlayer]triggered")

	member := message.Member
	// hack coz member.user in message.member is nil
	member.User = message.Author
	nickname := getNickname(member)
	data := database.PlayerData{GuildID: message.GuildID, UserID: message.Author.ID, Username: nickname}

	err := c.db.AddPlayer(ctx, &data)
	if err != nil {
		if err.Error() == database.PlayerAlreadyExistError {
			c.log.Debug().Err(err).Msgf("[commands.AddPlayer]player already exist")
			text := fmt.Sprintf("э, %s, куда тебе? Ты же уже участвуешь", data.Username)
			discord.ChannelMessageSend(message.ChannelID, text)
			return
		}
		c.log.Error().Err(err).Msgf("[commands.AddPlayer]Error. Can't add player")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	text := fmt.Sprintf("%s, вы приняты в пидорскую рулетку", data.Username)
	discord.ChannelMessageSend(message.ChannelID, text)
}

func (c *Commands) List(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Info().Msg("[commands.List]triggered")

	data, err := c.db.GetAllPlayers(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.List]Error. Can't get players")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	c.log.Debug().Msgf("[commands.List]Got %d players", len(data))

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
	c.log.Info().Msg("[commands.EventList]triggered")

	data, err := c.db.GetEvents(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.EventList]Error. Can't get events")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	c.log.Debug().Msgf("[commands.EventList]Got %d events", len(data))

	top := ""
	for i, d := range data {
		startTS := tools.ToDiscordTimeStamp(d.StartTime, tools.TSFormat().LongDateShortTime)
		endTS := tools.ToDiscordTimeStamp(d.EndTime, tools.TSFormat().LongDateShortTime)
		top += fmt.Sprintf("%d) %s. %s - %s\n", i+1, d.Type, startTS, endTS)
	}

	text := fmt.Sprintf("Список запланированных на этом сервере событий:\n%s", top)
	discord.ChannelMessageSend(message.ChannelID, text)
}

func (c *Commands) UpdatePlayersData(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Info().Msg("[commands.UpdatePlayersData]triggered")

	allPlayers, err := c.db.GetAllPlayers(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.UpdatePlayersData]Error. Can't get players")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	c.log.Debug().Msgf("[commands.UpdatePlayersData]Got %d players", len(allPlayers))

	text := fmt.Sprintf("Обновляю данные %d игроков...", len(allPlayers))
	discord.ChannelMessageSend(message.ChannelID, text)

	players := make([]*database.PlayerData, 0, len(allPlayers))
	for _, player := range allPlayers {
		u, err := discord.GuildMember(player.GuildID, player.UserID)
		if err != nil {
			c.log.Debug().Err(err).Msgf("[commands.UpdatePlayersData]Can't get player gid: %s; uid: %s; ", player.GuildID, player.UserID)
			continue
		}

		username := getNickname(u)
		if player.Username != username {
			p := &database.PlayerData{GuildID: player.GuildID, UserID: player.UserID, Score: player.Score, Username: username}
			players = append(players, p)
		}
	}

	err = c.db.UpdatePlayersData(ctx, players)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.UpdatePlayersData]Error. Can't update players data")
		discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		return
	}

	if len(players) == 1 {
		text = fmt.Sprintf("Данные %d игрока успешно обновлены", len(players))
	} else {
		text = fmt.Sprintf("Данные %d игроков успешно обновлены", len(players))
	}
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

	i := len(players)
	random := rand.New(rand.NewSource(time.Now().Unix()))
	num := random.Int31n(int32(i))

	result = players[num]

	return result, nil
}

func getNickname(member *discordgo.Member) string {
	nickname := member.Nick
	if nickname == "" {
		// TODO for discordgo pre 0.27.2, delete when package updates
		//nickname = cases.Title(language.Und, cases.NoLower).String(message.Author.Username)
		nickname = member.User.GlobalName
	}
	if nickname == "" {
		nickname = member.User.Username
	}

	return nickname
}
