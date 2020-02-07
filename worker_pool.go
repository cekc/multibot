package multibot

type WorkerPool interface {
	Submit(func())
	Wait()
}
