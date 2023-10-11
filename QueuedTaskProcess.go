package taskstore

import (
	"strings"
	"time"
)

func (store *Store) QueuedTaskProcess(queuedTask Queue) bool {
	// 1. Start queued task
	queuedTask.AppendDetails("Task started")
	attempts := queuedTask.Attempts + 1
	timeNow := time.Now()

	queuedTask.Status = QueueStatusRunning
	queuedTask.Attempts = attempts
	queuedTask.StartedAt = &timeNow
	store.QueueUpdate(&queuedTask)

	// 2. Find task definition
	task := store.TaskFindByID(queuedTask.TaskID)

	if task == nil {
		timeNow = time.Now()
		queuedTask.AppendDetails("Task DOES NOT exist")
		queuedTask.Status = QueueStatusFailed
		queuedTask.CompletedAt = &timeNow
		store.QueueUpdate(&queuedTask)
		return false
	}

	handlerFunc := store.taskHandlerFunc(task.Alias)

	result := handlerFunc(&queuedTask)

	if result {
		queuedTask.AppendDetails("Task completed")
		store.QueueSuccess(&queuedTask)
	} else {
		queuedTask.AppendDetails("Task failed")
		store.QueueFail(&queuedTask)
	}

	return true
}

// taskHandlerFunc finds the TaskHandler for the queued task and returns
// the Handle function, if not found, a default Handle function is returned
// which will print "No handler for alias ALIASNAME" message to notify the
// queue admin
func (store *Store) taskHandlerFunc(taskAlias string) func(queuedTask *Queue) bool {
	unifyName := func(name string) string {
		name = strings.ReplaceAll(name, "-", "")
		name = strings.ReplaceAll(name, "_", "")
		return name
	}

	for _, taskHandler := range store.taskHandlers {
		if strings.EqualFold(unifyName(taskHandler.Alias()), unifyName(taskAlias)) {
			return func(queuedTask *Queue) bool {
				taskHandler.SetQueuedTask(queuedTask)
				return taskHandler.Handle()
			}
		}
	}

	// for i := 0; i < len(store.taskHandlers); i++ {
	// 	if strings.EqualFold(unifyName(store.taskHandlers[i].Alias()), unifyName(taskAlias)) {
	// 		return func(queuedTask *Queue) bool {

	// 			return store.taskHandlers[i].Handle(TaskHandlerOptions{
	// 				QueuedTask: queuedTask,
	// 			})
	// 		}
	// 	}
	// }

	return func(queuedTask *Queue) bool {
		queuedTask.AppendDetails("No handler for alias: " + taskAlias)
		store.QueueUpdate(queuedTask)
		return false
	}
}
