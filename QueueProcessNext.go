package taskstore

import "log"

func (store *Store) QueueProcessNext() {
	runningTasks := store.QueueFindRunning(1)

	if len(runningTasks) > 0 {
		log.Println("There is already a running task " + runningTasks[0].ID + " (#" + runningTasks[0].ID + "). Queue stopped while completed'")
		return
	}

	nextQueuedTask := store.QueueFindNextQueuedTask()

	if nextQueuedTask == nil {
		// DEBUG log.Println("No queued tasks")
		return
	}

	store.QueuedTaskProcess(*nextQueuedTask)
}
