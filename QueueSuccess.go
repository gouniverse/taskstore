package taskstore

import "time"

// QueueSuccess completes a queued task  successfully
func (st *Store) QueueSuccess(queue *Queue) error {
	completedAt := time.Now()
	queue.CompletedAt = &completedAt
	queue.Status = QueueStatusSuccess
	return st.QueueUpdate(queue)
}
