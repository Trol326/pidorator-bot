package tools

import "fmt"

type timeStampFormat struct {
	ShortTime             string
	LongTime              string
	ShortDate             string
	LongDate              string
	LongDateShortTime     string
	VeryLongDateShortTime string
	Relative              string
}

func TSFormat() timeStampFormat {
	return timeStampFormat{
		ShortTime:             "t",
		LongTime:              "T",
		ShortDate:             "d",
		LongDate:              "D",
		LongDateShortTime:     "f",
		VeryLongDateShortTime: "F",
		Relative:              "R",
	}
}

func ToDiscordTimeStamp(time int64, format string) string {
	return fmt.Sprintf("<t:%d:%s>", time, format)
}
func ToDiscordTimeStamp32(time int32, format string) string {
	return fmt.Sprintf("<t:%d:%s>", time, format)
}

func UserIDToMention(UserID string) string {
	return fmt.Sprintf("<@%s>", UserID)
}
