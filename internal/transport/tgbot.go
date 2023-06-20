package transport

import (
	"os"
	"sync"

	// https://go-telegram-bot-api.dev/
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	botApi  *tgbotapi.BotAPI
	Wg      *sync.WaitGroup
	TimeOut int  // 30
	Debug   bool // false
	running bool
}

func (tgbot *TelegramBot) Run() {
	if !tgbot.running {
		tgbot.running = true
		defer func() { tgbot.running = false }()
		defer tgbot.Wg.Done()
		bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
		if err != nil {
			panic(err)
		}
		tgbot.botApi = bot
		bot.Debug = tgbot.Debug
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = tgbot.TimeOut
		updates := bot.GetUpdatesChan(updateConfig)
		for update := range updates {
			if update.Message != nil {
				// Запись запроса в БД
				// ...
				// Обратная связь: отправка ответа пользователю
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Принято в обработку, ожидайте ответа.")
				msg.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(msg); err != nil {
					panic(err) // ???
				}
			}
		}
		tgbot.botApi = nil
	}
}
