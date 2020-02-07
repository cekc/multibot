package multibot

import (
	"context"
)

type Handler interface {
	Handle(context.Context, Update)
}
