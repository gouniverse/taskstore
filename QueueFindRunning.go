package taskstore

func (store *Store) QueueFindRunning(limit int) []Queue {
	runningTasks, errList := store.QueueList(QueueListOptions{
		Status:    QueueStatusRunning,
		Limit:     limit,
		SortBy:    "created_at",
		SortOrder: "asc",
	})

	if errList != nil {
		return nil
	}

	return runningTasks
}
