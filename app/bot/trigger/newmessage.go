package trigger

import (
	"context"
	"fmt"
	"pidorator-bot/app/bot/command"
	"pidorator-bot/tools/roles"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	DefaultPrefix string = "!"
	GameEnabled   bool   = true
	AdminEnabled  bool   = true
)

/*
Trigger when user send new message on server. Only in available for bot channels

TODO maybe make trigger for game and for admin(?)
*/
func (t *Trigger) OnNewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}

	var prefix string
	ctx := context.Background()
	data, err := t.commands.GetBotData(ctx, message.GuildID)
	if err != nil {
		t.Log.Err(err).Msg("[trigger.OnNewMessage]error on getbotdata")
		prefix = DefaultPrefix
	} else {
		prefix = data.BotPrefix
	}

	if !strings.HasPrefix(message.Content, prefix) {
		return
	}

	t.Log.Info().Msgf("Got %s", message.Content)

	var game command.Game
	var admin command.Admin

	if GameEnabled && data.IsGameEnabled {
		game = t.commands
	}
	if AdminEnabled {
		admin = t.commands
	}

	switch {
	case strings.HasPrefix(message.Content, prefix+"help"):
		// TODO autodoc for help command
		_, err := discord.ChannelMessageSend(message.ChannelID, "Документация будет позже :D")
		if err != nil {
			t.Log.Err(err).Msg("[trigger.OnNewMessage]error on channelMessageSend")
			return
		}
	case strings.HasPrefix(message.Content, prefix+"ктопидор"):
		if game != nil {
			event, err := game.Who(ctx, discord, message)
			if err != nil {
				return
			}
			if data.IsAutoRollEnabled {
				t.OnEventCreation(ctx, discord, event)
			}
		}
	case strings.HasPrefix(message.Content, prefix+"списокпидоров") || strings.HasPrefix(message.Content, prefix+"топпидоров"):
		if game != nil {
			game.List(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, prefix+"disableautoroll") || strings.HasPrefix(message.Content, prefix+"changeautoroll"):
		if game != nil {
			game.ChangeAutoRoll(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, prefix+"пидордня"):
		if game != nil {
			game.AddPlayer(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, prefix+"обновитьпидоров"):
		if game != nil {
			game.UpdatePlayersData(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, prefix+"списоксобытий"):
		if game != nil {
			game.EventList(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, prefix+"botrename"):
		if admin != nil {
			admin.BotRename(ctx, discord, message)
		}
	case strings.HasPrefix(message.Content, prefix+"setprefix"):
		if admin != nil {
			// TODO refactor this later
			result, err := roles.MemberHasPermission(discord, message.GuildID, message.Author.ID, discordgo.PermissionAdministrator)
			if err != nil {
				_, errM := discord.ChannelMessageSend(message.ChannelID, "Sorry, server-side error. Please contact the bot maintainer")
				if errM != nil {
					t.Log.Err(errM).Msg("[trigger.OnNewMessage]error on channelMessageSend")
				}
				return
			}
			if !result {
				_, errM := discord.ChannelMessageSend(message.ChannelID, "У вас нет прав на использование этой команды")
				if errM != nil {
					t.Log.Err(errM).Msg("[trigger.OnNewMessage]error on channelMessageSend")
				}
				return
			}
			res, err := ParseCommand(prefix, message.Content)
			if err != nil {
				t.Log.Err(err).Msg("[trigger.OnNewMessage]error on message parsing")
				text := fmt.Sprintf("Ошибка при парсинге. Пример использования команды: `%ssetprefix новыйПрефикс`, например, `%ssetprefix ?`", prefix, prefix)
				_, errM := discord.ChannelMessageSend(message.ChannelID, text)
				if errM != nil {
					t.Log.Err(errM).Msg("[trigger.OnNewMessage]error on channelMessageSend")
				}
				return
			}
			if len(res) <= 2 {
				t.Log.Err(err).Msg("[trigger.OnNewMessage]not enough arguments")
				text := fmt.Sprintf("Ошибка при парсинге. Пример использования команды: `%ssetprefix новыйПрефикс`, например, `%ssetprefix ?`", prefix, prefix)
				_, errM := discord.ChannelMessageSend(message.ChannelID, text)
				if errM != nil {
					t.Log.Err(errM).Msg("[trigger.OnNewMessage]error on channelMessageSend")
				}
				return
			}
			t.Log.Info().Msgf("Parsed result: %s", res[2])
			admin.SetPrefix(ctx, discord, message, res[2])
		}
	default:
		t.Log.Debug().Msg("Command not found")
	}
}
