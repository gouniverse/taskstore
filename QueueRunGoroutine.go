package taskstore

import "time"

// QueueRunGoroutine goroutine to run the queue
//
// Example:
// go myTaskStore.QueueRunGoroutine(10, 2)
//
// Params:
// - processSeconds int - time to wait until processing the next task (i.e. 10s)
// - unstuckMinutes int - time to wait before mark running tasks as failed
func (store *Store) QueueRunGoroutine(processSeconds int, unstuckMinutes int) {
	i := 0
	for {
		i++

		store.QueueUnstuck(unstuckMinutes)

		time.Sleep(1 * time.Second) // Sleep 1 second

		store.QueueProcessNext()

		time.Sleep(time.Duration(processSeconds) * time.Second) // Every 10 seconds
	}
}
