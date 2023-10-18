package taskstore

import "time"

func (store *Store) QueuedTaskForceFail(queuedTask Queue, waitMinutes int) error {
	startedAt := queuedTask.StartedAt

	if startedAt == nil {
		return nil
	}

	minutes := -1 * waitMinutes

	if startedAt.Before(time.Now().Add((time.Duration(minutes) * time.Minute))) {
		queuedTask.AppendDetails("Failed forcefully after 2 minutes timeout")
		completedAt := time.Now()
		queuedTask.CompletedAt = &completedAt
		queuedTask.Status = QueueStatusFailed
		return store.QueueUpdate(&queuedTask)
	}

	return nil
}
