package taskstore

func (store *Store) TaskHandlerList() []TaskHandlerInterface {
	return store.taskHandlers
}
