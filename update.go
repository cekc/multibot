package multibot

type Update interface {
	Body() string
	Hook() Notifier
}
