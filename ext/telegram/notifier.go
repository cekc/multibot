package telegram

import (
	"context"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Notifier struct {
	bot        *tgbotapi.BotAPI
	telegramID int
}

func (notifier Notifier) Notify(ctx context.Context, message string) {
	msg := tgbotapi.NewMessage(int64(notifier.telegramID), message)
	_, err := notifier.bot.Send(msg) // log errors here!

	if err != nil {
		log.Println("Notify() failed:", err)
	}
}
