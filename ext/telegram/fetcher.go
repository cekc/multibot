package telegram

import (
	"context"
	"net/url"
	"path"
	"time"

	"github.com/cekc/multibot"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

type Fetcher struct {
	bot        *tgbotapi.BotAPI
	rawUpdates tgbotapi.UpdatesChannel
	debug      bool
}

func NewFetcher() *Fetcher {
	return &Fetcher{}
}

func (fetcher *Fetcher) WithDebug(debug bool) *Fetcher {
	fetcher.debug = debug
	return fetcher
}

func (fetcher *Fetcher) Webhook(token, webhookURLPrefix string, certFile interface{}) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return errors.WithMessage(err, "Could not connect to Telegam bot API")
	}

	webhookURL := joinURL(webhookURLPrefix, token)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert(webhookURL, certFile))
	if err != nil {
		return errors.WithMessage(err, "Could not set webhook for Telegram bot")
	}

	bot.Debug = fetcher.debug
	fetcher.bot = bot
	fetcher.rawUpdates = bot.ListenForWebhook(extractURLPath(webhookURL))

	return nil
}

func (fetcher *Fetcher) LongPoll(token string, offset int, timeout time.Duration) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return errors.WithMessage(err, "Could not connect to Telegam bot API")
	}

	updateConfig := tgbotapi.NewUpdate(offset)
	updateConfig.Timeout = int(timeout.Seconds() + 0.5)

	rawUpdates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		return errors.WithMessage(err, "Could not start polling Telegram bot API")
	}

	bot.Debug = fetcher.debug
	fetcher.bot = bot
	fetcher.rawUpdates = rawUpdates

	return nil
}

func (fetcher *Fetcher) Fetch(ctx context.Context) <-chan multibot.Update {
	updates := make(chan multibot.Update)

	go func() {
		defer close(updates)

		for {
			select {
			case <-ctx.Done():
				return

			case rawUpdate, ok := <-fetcher.rawUpdates:
				if !ok {
					return
				}
				if rawUpdate.Message != nil {
					updates <- Update{
						Text: rawUpdate.Message.Text,
						Notifier: Notifier{
							bot:        fetcher.bot,
							telegramID: rawUpdate.Message.From.ID,
						},
					}
				}
			}
		}
	}()

	return updates
}

func joinURL(prefix, suffix string) string {
	u, err := url.Parse(prefix)
	if err != nil {
		panic(err)
	}

	u.Path = path.Join(u.Path, suffix)
	return u.String()
}

func extractURLPath(addr string) string {
	u, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}

	return u.Path
}
