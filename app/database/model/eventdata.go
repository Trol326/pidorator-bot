package model

import (
	"fmt"
	"time"
)

type EventData struct {
	GuildID   string `bson:"guildID,omitempty"`
	ChannelID string `bson:"channelID,omitempty"`
	Type      string `bson:"eventType,omitempty"`
	StartTime int64  `bson:"startTime"`
	EndTime   int64  `bson:"endTime"`
}

func (d *EventData) String() string {
	return fmt.Sprintf("{GuildID: %s, ChannelID: %s, Type: %s, StartTime: %d, EndTime: %d}", d.GuildID, d.ChannelID, d.Type, d.StartTime, d.EndTime)
}

func (d EventData) IsEventEnded(now ...int64) bool {
	if d.StartTime > d.EndTime {
		return true
	}
	if d.StartTime == 0 || d.EndTime == 0 {
		return false
	}
	if len(now) < 1 {
		return false
	}
	return d.EndTime < now[0]
}

func (d EventData) SecondsUntilEnd() int64 {
	var result int64 = 0
	now := time.Now().Unix()
	if d.IsEventEnded(now) {
		return result
	}
	result = d.EndTime - now

	return result
}
