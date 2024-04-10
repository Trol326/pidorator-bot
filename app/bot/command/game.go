package command

import (
	"context"
	"fmt"
	"math/rand"
	"pidorator-bot/app/database"
	"pidorator-bot/app/database/model"
	"pidorator-bot/tools"
	"pidorator-bot/tools/roles"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (c *Commands) AutoRoll(ctx context.Context, discord *discordgo.Session, guildID string, channelID string) (*model.EventData, error) {
	c.log.Info().Msg("[commands.Autoroll]triggered")

	data, err := c.db.GetBotData(ctx, guildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.DisableAutoRoll]Error. Can't get bot data")
		_, errM := discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.AutoRoll]error on channelMessageSend")
		}
		return nil, err
	}

	if !data.IsAutoRollEnabled {
		_, errM := discord.ChannelMessageSend(channelID, "Сорян, накосячил. Автокрутки отключены")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.AutoRoll]error on channelMessageSend")
		}
		err := fmt.Errorf("autorolls disabled")
		return nil, err
	}

	ev, err := c.db.FindEvent(ctx, guildID, database.GameEventName)
	if err != nil && err.Error() != database.EventAlreadyExistError {
		c.log.Error().Err(err).Msgf("[commands.Autoroll]Error. Can't find event")
		_, errM := discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.AutoRoll]error on channelMessageSend")
		}
		return nil, err
	}
	if ev.GuildID != "" && !IsGameEventEnded(ev) {
		c.log.Debug().Msg("[commands.Autoroll]Can't roll, event still not ended")
		text := fmt.Sprintf("Накосячил, ещё рано. В предыдущий раз крутили %s", tools.ToDiscordTimeStamp(ev.StartTime, tools.TSFormat().Relative))
		_, errM := discord.ChannelMessageSend(channelID, text)
		if errM != nil {
			c.log.Err(errM).Msg("[commands.AutoRoll]error on channelMessageSend")
		}
		err := fmt.Errorf("event not ended")
		return nil, err
	}

	player, err := c.getRandomPlayer(ctx, guildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.Autoroll]Error. Can't get player")
		_, errM := discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.AutoRoll]error on channelMessageSend")
		}
		return nil, err
	}

	err = c.db.IncreasePlayerScore(ctx, player.GuildID, player.UserID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.Autoroll]Error. Can't increase player score")
		_, errM := discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.AutoRoll]error on channelMessageSend")
		}
		return nil, err
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Hour * 24)
	event := model.EventData{
		GuildID:   guildID,
		ChannelID: channelID,
		Type:      database.GameEventName,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
	}
	err = c.db.AddEvent(ctx, &event)
	if err != nil && err.Error() != database.EventAlreadyExistError {
		c.log.Error().Err(err).Msgf("[commands.Autoroll]Error. Can't add event")
		_, errM := discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.AutoRoll]error on channelMessageSend")
		}
		return nil, err
	}
	if err.Error() == database.EventAlreadyExistError {
		err := c.db.UpdateEvent(ctx, &event)
		if err != nil {
			c.log.Error().Err(err).Msgf("[commands.Autoroll]Error. Can't update event")
			_, errM := discord.ChannelMessageSend(channelID, "Sorry, server-side error. Please contact the bot maintainer")
			if errM != nil {
				c.log.Err(errM).Msg("[commands.AutoRoll]error on channelMessageSend")
			}
			return nil, err
		}
	}

	text := fmt.Sprintf("%s, ты пидор <:MumeiYou:1192139708222935050>", tools.UserIDToMention(player.UserID))
	_, errM := discord.ChannelMessageSend(channelID, text)
	if errM != nil {
		c.log.Err(errM).Msg("[commands.AutoRoll]error on channelMessageSend")
	}

	err = c.finishRoll(ctx, discord, data, player)
	if err != nil {
		c.log.Err(err).Msg("[commands.Who]error. Can't finishRoll")
		_, errM := discord.ChannelMessageSend(channelID, "Ошибка назначения роли пидора дня")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.Who]error on channelMessageSend")
		}
	}

	return &event, nil
}

func (c *Commands) Who(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) (*model.EventData, error) {
	c.log.Info().Msg("[commands.Who]triggered")

	data, err := c.db.GetBotData(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't get bot data")
		_, errM := discord.ChannelMessageSend(message.GuildID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.Who]error on channelMessageSend")
		}
		return nil, err
	}

	if !data.IsGameEnabled {
		_, errM := discord.ChannelMessageSend(message.GuildID, "Игра отключена")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.Who]error on channelMessageSend")
		}
		err := fmt.Errorf("game disabled")
		return nil, err
	}

	ev, err := c.db.FindEvent(ctx, message.GuildID, database.GameEventName)
	if err != nil && err.Error() != database.EventAlreadyExistError && err.Error() != c.db.NotFoundError() {
		c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't find event")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(err).Msg("[commands.Who]error on channelMessageSend")
			return nil, errM
		}
		return nil, err
	}
	if ev != nil && ev.GuildID != "" && !IsGameEventEnded(ev) {
		c.log.Debug().Msg("[commands.Who]Can't roll, event still not ended")
		text := fmt.Sprintf("Да ты погоди, %s, время ещё не пришло. В предыдущий раз крутили %s", message.Author.GlobalName, tools.ToDiscordTimeStamp(ev.StartTime, tools.TSFormat().Relative))
		_, errM := discord.ChannelMessageSend(message.ChannelID, text)
		if errM != nil {
			c.log.Err(errM).Msg("[commands.Who]error on channelMessageSend")
		}
		err = fmt.Errorf("event not ended")
		return nil, err
	}

	player, err := c.getRandomPlayer(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't get player")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.Who]error on channelMessageSend")
		}
		return nil, err
	}

	err = c.db.IncreasePlayerScore(ctx, player.GuildID, player.UserID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't increase player score")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.Who]error on channelMessageSend")
		}
		return nil, err
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Hour * 24)
	event := model.EventData{
		GuildID:   message.GuildID,
		ChannelID: message.ChannelID,
		Type:      database.GameEventName,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
	}
	err = c.db.AddEvent(ctx, &event)
	if err != nil && err.Error() != database.EventAlreadyExistError {
		c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't add event")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.Who]error on channelMessageSend")
		}
		return nil, err
	}
	if err != nil && err.Error() == database.EventAlreadyExistError {
		err := c.db.UpdateEvent(ctx, &event)
		if err != nil {
			c.log.Error().Err(err).Msgf("[commands.Who]Error. Can't update event")
			_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
			if errM != nil {
				c.log.Err(errM).Msg("[commands.Who]error on channelMessageSend")
			}
			return nil, err
		}
	}

	text := fmt.Sprintf("%s, ты пидор <:MumeiYou:1192139708222935050>", tools.UserIDToMention(player.UserID))
	_, errM := discord.ChannelMessageSend(message.ChannelID, text)
	if errM != nil {
		c.log.Err(errM).Msg("[commands.Who]error on channelMessageSend")
	}

	err = c.finishRoll(ctx, discord, data, player)
	if err != nil {
		c.log.Err(err).Msg("[commands.Who]error. Can't finishRoll")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Ошибка назначения роли пидора дня")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.Who]error on channelMessageSend")
		}
	}

	return &event, nil
}

func (c *Commands) ChangeAutoRoll(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Info().Msg("[commands.DisableAutoRoll]triggered")

	data, err := c.db.GetBotData(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.DisableAutoRoll]Error. Can't get bot data")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(err).Msg("[commands.DisableAutoRoll]error on channelMessageSend")
			return
		}
		return
	}

	data.IsAutoRollEnabled = !data.IsAutoRollEnabled
	err = c.db.ChangeBotData(ctx, data)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.DisableAutoRoll]Error. Can't change bot data")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(err).Msg("[commands.DisableAutoRoll]error on channelMessageSend")
			return
		}
		return
	}

	status := ""
	if data.IsAutoRollEnabled {
		status = "включены"
	} else {
		status = "отключены"
	}

	text := fmt.Sprintf("Теперь автокрутки %s", status)
	_, errM := discord.ChannelMessageSend(message.ChannelID, text)
	if errM != nil {
		c.log.Err(err).Msg("[commands.DisableAutoRoll]error on channelMessageSend")
		return
	}
}

func (c *Commands) AddPlayer(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Info().Msg("[commands.AddPlayer]triggered")

	member := message.Member
	// hack coz member.user in message.member is nil
	member.User = message.Author
	nickname := getNickname(member)
	data := model.PlayerData{GuildID: message.GuildID, UserID: message.Author.ID, Username: nickname}

	err := c.db.AddPlayer(ctx, &data)
	if err != nil {
		if err.Error() == database.PlayerAlreadyExistError {
			c.log.Debug().Err(err).Msgf("[commands.AddPlayer]player already exist")
			text := fmt.Sprintf("э, %s, куда тебе? Ты же уже участвуешь", data.Username)
			_, errM := discord.ChannelMessageSend(message.ChannelID, text)
			if errM != nil {
				c.log.Err(errM).Msg("[commands.AddPlayer]error on channelMessageSend")
			}
			return
		}
		c.log.Error().Err(err).Msgf("[commands.AddPlayer]Error. Can't add player")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.AddPlayer]error on channelMessageSend")
		}
		return
	}

	text := fmt.Sprintf("%s, вы приняты в пидорскую рулетку", data.Username)
	_, errM := discord.ChannelMessageSend(message.ChannelID, text)
	if errM != nil {
		c.log.Err(errM).Msg("[commands.AddPlayer]error on channelMessageSend")
	}
}

func (c *Commands) List(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Info().Msg("[commands.List]triggered")

	data, err := c.db.GetAllPlayers(ctx, message.GuildID, database.DcendingSorting)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.List]Error. Can't get players")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.List]error on channelMessageSend")
		}
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
	_, errM := discord.ChannelMessageSend(message.ChannelID, text)
	if errM != nil {
		c.log.Err(errM).Msg("[commands.List]error on channelMessageSend")
	}
}

func (c *Commands) EventList(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Info().Msg("[commands.EventList]triggered")

	data, err := c.db.GetEvents(ctx, message.GuildID)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.EventList]Error. Can't get events")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.EventList]error on channelMessageSend")
		}
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
	_, errM := discord.ChannelMessageSend(message.ChannelID, text)
	if errM != nil {
		c.log.Err(errM).Msg("[commands.EventList]error on channelMessageSend")
	}
}

func (c *Commands) UpdatePlayersData(ctx context.Context, discord *discordgo.Session, message *discordgo.MessageCreate) {
	c.log.Info().Msg("[commands.UpdatePlayersData]triggered")

	allPlayers, err := c.db.GetAllPlayers(ctx, message.GuildID, database.NoSorting)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.UpdatePlayersData]Error. Can't get players")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.UpdatePlayersData]error on channelMessageSend")
		}
		return
	}

	c.log.Debug().Msgf("[commands.UpdatePlayersData]Got %d players", len(allPlayers))

	text := fmt.Sprintf("Обновляю данные %d игроков...", len(allPlayers))
	_, errM := discord.ChannelMessageSend(message.ChannelID, text)
	if errM != nil {
		c.log.Err(errM).Msg("[commands.UpdatePlayersData]error on channelMessageSend")
	}

	players := make([]*model.PlayerData, 0, len(allPlayers))
	for _, player := range allPlayers {
		u, err := discord.GuildMember(player.GuildID, player.UserID)
		if err != nil {
			c.log.Debug().Err(err).Msgf("[commands.UpdatePlayersData]Can't get player gid: %s; uid: %s; ", player.GuildID, player.UserID)
			continue
		}

		username := getNickname(u)
		if player.Username != username {
			p := &model.PlayerData{GuildID: player.GuildID, UserID: player.UserID, Score: player.Score, Username: username}
			players = append(players, p)
		}
	}

	err = c.db.UpdatePlayersData(ctx, players)
	if err != nil {
		c.log.Error().Err(err).Msgf("[commands.UpdatePlayersData]Error. Can't update players data")
		_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
		if errM != nil {
			c.log.Err(errM).Msg("[commands.UpdatePlayersData]error on channelMessageSend")
		}
		return
	}

	if len(players) == 1 {
		text = fmt.Sprintf("Данные %d игрока успешно обновлены", len(players))
	} else {
		text = fmt.Sprintf("Данные %d игроков успешно обновлены", len(players))
	}
	_, errM = discord.ChannelMessageSend(message.ChannelID, text)
	if errM != nil {
		c.log.Err(errM).Msg("[commands.UpdatePlayersData]error on channelMessageSend")
	}
}

// Checks is game event ended by time. True if ended, False otherwise
func IsGameEventEnded(event *model.EventData) bool {
	if event == nil {
		return true
	}
	now := time.Now()
	res := event.IsEventEnded(now.Unix())
	if res {
		return res
	}

	endTime := time.Unix(event.StartTime, 0).Add(time.Hour * 24)

	// if startTime + 24h < now => returns true
	return endTime.Before(now)
}

func (c *Commands) getRandomPlayer(ctx context.Context, guildID string) (*model.PlayerData, error) {
	result := &model.PlayerData{}

	players, err := c.db.GetAllPlayers(ctx, guildID, database.NoSorting)
	if err != nil {
		return result, err
	}

	i := len(players)
	random := rand.New(rand.NewSource(time.Now().Unix()))
	num := random.Int31n(int32(i))

	result = players[num]

	return result, nil
}

func (c *Commands) finishRoll(ctx context.Context, discord *discordgo.Session, data *model.BotData, winner *model.PlayerData) error {
	newData := *data

	if data.LastPidorRoleID == "" {
		c.log.Debug().Msg("lastPidorRoleID not found")
		role, err := roles.CreateRole(discord, newData.GuildID)
		if err != nil {
			c.log.Err(err).Msg("[commands.finishRoll]error on role creation")
			return err
		}
		newData.LastPidorRoleID = role.ID
	}

	r, err := roles.GetRole(discord, newData.GuildID, newData.LastPidorRoleID)
	if err != nil || r.ID == "" {
		c.log.Err(err).Msg("[commands.finishRoll]error. Can't get role")
		return err
	}

	if data.LastPidorUserID != "" && data.LastPidorRoleID != "" {
		err := roles.DeleteUserRole(discord, data.GuildID, data.LastPidorUserID, data.LastPidorRoleID)
		if err != nil {
			c.log.Err(err).Msg("[commands.finishRoll]error on delete old winner role")
			return err
		}
	}

	err = roles.SetUserRole(discord, newData.GuildID, winner.UserID, newData.LastPidorRoleID)
	if err != nil {
		c.log.Err(err).Msg("[commands.finishRoll]error on SetUserRole")
		return err
	}
	newData.LastPidorUserID = winner.UserID

	err = c.db.ChangeBotData(ctx, &newData)
	if err != nil {
		c.log.Err(err).Msg("[commands.finishRoll]error on changeBotData")
		return err
	}

	return nil
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
