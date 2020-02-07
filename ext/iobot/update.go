package iobot

import (
	"github.com/cekc/multibot"
)

type Update struct {
	Text     string
	Notifier *Notifier
}

func (update Update) Body() string {
	return update.Text
}

func (update Update) Hook() multibot.Notifier {
	return update.Notifier
}
