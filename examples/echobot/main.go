package main

import (
	"context"
	"fmt"

	"github.com/cekc/multibot"
	"github.com/cekc/multibot/ext/filebot"
)

type echoer struct {
	name string
}

func (e echoer) Handle(ctx context.Context, update multibot.Update) {
	fmt.Println(update.From, e.name, "says:", update.Body)
}

func main() {
	bot := multibot.New()

	bot.RegisterFetchers(filebot.NewConsoleFetcher())
	bot.RegisterHandlers(echoer{"Mike"}, echoer{"Zuloos"})

	bot.Serve()

	fmt.Println("Gracefully shut down")
}
