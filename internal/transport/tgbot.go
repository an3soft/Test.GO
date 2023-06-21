package transport

import (
	c "an3softbot/internal/contracts"
	m "an3softbot/internal/models"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	// https://go-telegram-bot-api.dev/
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Writer  c.DBWriter
	Wg      *sync.WaitGroup
	TimeOut int  // 60
	Debug   bool // false
	botApi  *tgbotapi.BotAPI
	running bool
}

func (tgbot *TelegramBot) Init() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	tgbot.botApi = bot
	bot.Debug = tgbot.Debug
}

func (tgbot *TelegramBot) Run(ctx context.Context) {
	if !tgbot.running {
		tgbot.running = true
		defer func() { tgbot.running = false }()
		defer tgbot.Wg.Done()
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = tgbot.TimeOut
		bot := tgbot.botApi
		updates := bot.GetUpdatesChan(updateConfig)
		for {
			select {
			case update := <-updates:
				if update.Message != nil {
					if update.Message.Text == "/start" {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Добрый день, <b>%s</b>!", update.Message.From.UserName))
						msg.ParseMode = "html"
						if _, err := bot.Send(msg); err != nil {
							panic(err)
						}
					} else {
						// Регистрация запроса
						req := m.Request{}
						req.UserId = update.Message.From.ID
						req.ChatId = update.Message.Chat.ID
						req.MessageID = update.Message.MessageID
						req.UserName = update.Message.From.UserName
						req.Text = update.Message.Text
						req.Received = update.Message.Time()
						req.Updated = time.Now()
						tgbot.Writer.Write(ctx, &req)
						// Обратная связь: отправка ответа пользователю
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Принято в обработку, ожидайте ответа.")
						msg.ReplyToMessageID = update.Message.MessageID
						msg.ParseMode = "html"
						if _, err := bot.Send(msg); err != nil {
							panic(err)
						}
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func (tgbot *TelegramBot) Send(ctx context.Context, answer *m.Answer) {
	msg := tgbotapi.NewMessage(answer.ChatId, answer.Text)
	msg.ParseMode = "html"
	bot := tgbot.botApi
	if bot != nil {
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}
