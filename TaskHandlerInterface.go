package taskstore

type TaskHandlerInterface interface {
	Alias() string

	Title() string

	Description() string

	Handle() bool

	SetQueuedTask(queuedTask *Queue)

	SetOptions(options map[string]string)
}
