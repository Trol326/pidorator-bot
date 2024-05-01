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
			"Внимание, у нас пидор дня! Возможно, криминал! По коням!",
			"Woop-woop! That's the sound of da pidor-police!",
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
		{
			"Осторожно! Пидор дня активирован!",
			"Военный спутник запущен, коды доступа внутри...",
			"Что с нами стало...",
		},
		{
			"Кто рано встает, тому Pidorator звание пидора дня дает!",
			"Дзынь дзынь",
		},
		{
			"Надпись гласит: «Не влезай, мол, убьёт».",
			"Этот всё сделает наоборот,",
			"Влезет, дурак, и проблему создаст.",
			"Кто после этого он?",
		},
		{
			"Было лето, или утро, или тучи, или день,",
			"Или ветер, или вечер, или дождик, или тень.",
			"Было рано, или осень, или месяц, или пыль,",
			"То ли завтра, то ли в полдень...",
		},
		{
			"— Кто это?",
			"— Наверное, король...",
			"— Нет, это пидор дня!",
		},
		{
			"Осторожно! Пидор дня активирован!",
			"Выезжаю на место...",
			"Не может быть!",
		},
		{
			"Система понятна... на схемах отработана",
			"Вот она, вот она на пидоре намотана...",
		},
		{
			"А теперь нечто совсем иное... а нет, все то же самое...",
			"Определяем пидора дня методом тыка...",
			"Тык-тык... тык... тык...",
			"Натыкал пидора дня!",
		},
		{
			"Ведётся поиск в базе данных",
			"Так-так, что же тут у нас...",
			"Анализ завершен.",
		},
		{
			"Система взломана. Нанесён урон. Запущено планирование контрмер.",
			"Выезжаю на место...",
			"Ох...",
		},
		{
			"На дворе трава, на траве дрова",
			"А мы определяем пидора дня!",
			"<:pepelook:1235195238960857169>",
		},
		{
			"Новое масло «Whizzo» совершенно невозможно отличить от мертвого краба!",
			"А пидора дня от человека я отличить могу!",
			"<:pepelook:1235195238960857169>",
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
		"МЫ ОПРЕДЕЛИЛИ, ЧТО %s, ПИДОР!",
		"*Достает грамоту* Объявляю тебя, %s, пидором дня. Гип гип, ура!",
		"Ты, %s, пидор!",
		"Выборы, выборы в претенденты в пидоры!\n%s, ты самый главный этого дня",
		"Ученые строят гипотезы, как из живых существ может сформироваться такой пидор как, %s",
		"Pidorator рапортовал, да не дорапортовал, а стал дорапортовывать, зарапортовался и нарапортавал нам, пидора %s",
		"Один кинул – не докинул; другой кинул перекинул; третий кинул – не попал; четвертый — %s — пидором дня стал. ",
		"Инцидент с интендантом, прецедент с претендентом, интрига с интриганом, а @Hylegan с титулом пидора дня.",
		"Ого, вы посмотрите только! А пидор дня то — %s",
		"*Кряхтит* %s — ты пидор дня!",
		"Наши специалисты используют квантово-механические свойства конденсата Бозе-Эйнштейна для определения пидора дня.\nБлагодаря сложным вычислительным манипуляциям им удалось выявить, %s Эврика!",
	}

	num := tools.GetRandomInt32(len(resultPhrases))

	return resultPhrases[num]
}
