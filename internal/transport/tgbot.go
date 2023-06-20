package transport

import (
	c "an3softbot/internal/contracts"
	m "an3softbot/internal/models"
	"os"
	"sync"
	"time"

	// https://go-telegram-bot-api.dev/
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Wg      *sync.WaitGroup
	Writer  c.Writer
	TimeOut int  // 30
	Debug   bool // false
	botApi  *tgbotapi.BotAPI
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
			if update.Message != nil && update.Message.Text != "/start" {
				// Регистрация запроса
				req := m.Request{}
				req.UserId = update.Message.From.ID
				req.ChatId = update.Message.Chat.ID
				req.MessageID = update.Message.MessageID
				req.UserName = update.Message.From.UserName
				req.Text = update.Message.Text
				req.Received = update.Message.Time()
				req.Updated = time.Now()
				tgbot.Writer.Write(req)
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
