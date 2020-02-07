package multibot

import (
	"context"
)

type Fetcher interface {
	Fetch(context.Context) <-chan Update
}
