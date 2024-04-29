package content

import (
	"pidorator-bot/tools"
)

func GetRandomTeasePhrases() []string {
	var teasePhrases = [][]string{
		{
			"Woob-woob, that's da sound of da pidor-police!",
			"Выезжаю на место...",
			"Но кто же он?",
		},
		{
			"Woob-woob, that's da sound of da pidor-police!",
			"Ведётся поиск в базе данных",
			"Ведётся захват подозреваемого...",
		},
		{
			"Что тут у нас?",
			"А могли бы на работе делом заниматься...",
			"Проверяю данные...",
		},
		{
			"Инициирую поиск пидора дня...",
			"Машины выехали",
			"Так-так, что же тут у нас...",
		},
		{
			"Что тут у нас?",
			"Военный спутник запущен, коды доступа внутри...",
			"Не может быть!",
		},
	}

	num := tools.GetRandomInt32(len(teasePhrases))

	return teasePhrases[num]
}

func GetRandomResultPhrase() string {
	var resultPhrases = []string{
		"А вот и пидор - %s",
		"Вот ты и пидор, %s",
		"Ну ты и пидор, %s",
		"Сегодня ты пидор, %s",
		"Анализ завершен, сегодня ты пидор, %s",
		"ВЖУХ И ТЫ ПИДОР, %s",
		"Пидор дня обыкновенный, 1шт. - %s",
		"Стоять! Не двигаться! Вы объявлены пидором дня, %s",
		"И прекрасный человек дня сегодня... а нет, ошибка, всего-лишь пидор - %s",
		"%s, ты пидор <:MumeiYou:1192139708222935050>",
	}

	num := tools.GetRandomInt32(len(resultPhrases))

	return resultPhrases[num]
}
