package taskstore

type TaskHandlerInterface interface {
	Alias() string

	Title() string

	Description() string

	Handle(opts TaskHandlerOptions) bool
}
