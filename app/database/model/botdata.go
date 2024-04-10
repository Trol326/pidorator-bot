package model

import "fmt"

type BotData struct {
	GuildID           string `bson:"guildID"`
	GameChannelID     string `bson:"gameChannelID,omitempty"`
	BotPrefix         string `bson:"botPrefix,omitempty"`
	IsGameEnabled     bool   `bson:"isGameEnabled"`
	IsAutoRollEnabled bool   `bson:"isAutoRollEnabled"`
}

func (d *BotData) String() string {
	return fmt.Sprintf("{GuildID: %s, GameChannelID: %s, BotPrefix: %s}", d.GuildID, d.GameChannelID, d.BotPrefix)
}
