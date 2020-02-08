[![](https://godoc.org/github.com/cekc/multibot?status.svg)](http://godoc.org/github.com/cekc/multibot)

# multibot

## Types

### [`Multibot`](https://godoc.org/github.com/cekc/multibot#Multibot)

`Multibot` is nothing more than a collections of `Fetcher`s and a collection of `Handler`s.

`Fetcher`s and `Handler`s have different repsonsibilities: `Fetcher` is for communicating with API, `Handler` is for your business logic.

### [`Fetcher`](https://godoc.org/github.com/cekc/multibot#Fetcher)

`Fetcher` is just an interface whose single method fetches updates from some source and returns them in a `<-chan Update`. Plain simple.

But there is an important detail. Look more closely on `Update` type: it has method `Hook()` of return type `Notifier`

```go
type Notifier interface {
	Notify(ctx context.Context, message string)
}
```

This provides a shortcut for replying to a hypothetical user which has produced an update. The responsibilty of prividing right `Notifier` falls on `Fetcher`'s shoulders.