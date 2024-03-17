package taskstore

func (store *Store) QueueFindRunning(limit int) []Queue {
	runningTasks, errList := store.QueueList(QueueQueryOptions{
		Status:    QueueStatusRunning,
		Limit:     limit,
		SortBy:    COLUMN_CREATED_AT,
		SortOrder: ASC,
	})

	if errList != nil {
		return nil
	}

	return runningTasks
}
