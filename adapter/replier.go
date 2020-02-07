package adapter

import (
	"context"

	"github.com/cekc/multibot"
)

type Replier struct {
	Processor func(string) string
}

func (replier Replier) Handle(ctx context.Context, update multibot.Update) {
	update.Hook().Notify(ctx, replier.Processor(update.Body()))
}
