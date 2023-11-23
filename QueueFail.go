package taskstore

import "time"

// QueueFail fails a queued task
func (st *Store) QueueFail(queue *Queue) error {
	completedAt := time.Now()
	queue.CompletedAt = &completedAt
	queue.Status = QueueStatusFailed
	return st.QueueUpdate(queue)
}
