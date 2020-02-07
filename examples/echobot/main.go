package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cekc/multibot"
	"github.com/cekc/multibot/adapter"
	"github.com/cekc/multibot/ext/iobot"
	"github.com/cekc/multibot/shutdown"
)

func main() {
	bot := multibot.New()

	echoer1 := adapter.Replier{func(message string) string { return fmt.Sprint("Mike says: ", message) }}
	echoer2 := adapter.Replier{func(message string) string { return fmt.Sprint("Zu utters: ", message) }}

	bot.AddFetcher(iobot.NewConsoleFetcher())
	bot.AddHandler(echoer1)
	bot.AddHandler(echoer2)

	go bot.Process()

	select {
	case signal := <-shutdown.QuitSignal():
		log.Printf("Received signal <%v>, shutting down...\n", signal)
	case <-bot.RanOutOfUpdates():
		log.Println("No more updates, shutting down...")
	}

	err := bot.Shutdown(context.Background())

	if err == nil {
		log.Println("Gracefully shut down")
	} else {
		log.Println(err.Error())
	}
}
