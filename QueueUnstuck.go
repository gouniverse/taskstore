package taskstore

// QueueUnstuck clears the queue of tasks running for more than the
// specified wait time as most probably these have abnormally
// exited (panicked) and stop the rest of the queue from being
// processed
//
// The tasks are marked as failed. However, if they are still running
// in the background and they are successfully completed, they will
// be marked as success
//
// =================================================================
// Business Logic
// 1. Checks is there are running tasks in progress
// 2. If running for more than the specified wait minutes mark as failed
// =================================================================
func (store *Store) QueueUnstuck(waitMinutes int) {
	runningTasks := store.QueueFindRunning(3)

	if len(runningTasks) < 1 {
		return
	}

	for _, runningTask := range runningTasks {
		store.QueuedTaskForceFail(runningTask, waitMinutes)
	}
}
