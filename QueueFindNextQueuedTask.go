package taskstore

import "log"

// QueueFindNextQueuedTask finds the next queued task
// that is ready to be processed
func (store *Store) QueueFindNextQueuedTask() *Queue {
	queuedTasks, errList := store.QueueList(QueueQueryOptions{
		Status:    QueueStatusQueued,
		Limit:     1,
		SortBy:    COLUMN_CREATED_AT,
		SortOrder: ASC,
	})

	if errList != nil {
		log.Println(errList)
		return nil
	}

	if len(queuedTasks) < 1 {
		return nil
	}

	return &queuedTasks[0]
}
