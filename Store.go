package taskstore

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/utils"
	"github.com/mingrammer/cfmt"
)

// Store defines a session store
type Store struct {
	taskTableName      string
	queueTableName     string
	taskHandlers       []TaskHandlerInterface
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
}

var _ StoreInterface = (*Store)(nil)

// NewStoreOptions define the options for creating a new task store
type NewStoreOptions struct {
	TaskTableName      string
	QueueTableName     string
	DB                 *sql.DB
	DbDriverName       string
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new task store
func NewStore(opts NewStoreOptions) (*Store, error) {
	store := &Store{
		taskTableName:      opts.TaskTableName,
		queueTableName:     opts.QueueTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
	}

	if store.taskTableName == "" {
		return nil, errors.New("task store: taskTableName is required")
	}

	if store.queueTableName == "" {
		return nil, errors.New("task store: queueTableName is required")
	}

	if store.db == nil {
		return nil, errors.New("task store: DB is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = sb.DatabaseDriverName(store.db)
	}

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}

// AutoMigrate migrates the tables
func (st *Store) AutoMigrate() error {
	sqlTaskTable := st.SqlCreateTaskTable()

	if st.debugEnabled {
		log.Println(sqlTaskTable)
	}

	_, errTask := st.db.Exec(sqlTaskTable)
	if errTask != nil {
		log.Println(errTask)
		return errTask
	}

	sqlQueueTable := st.SqlCreateQueueTable()

	if st.debugEnabled {
		log.Println(sqlQueueTable)
	}

	_, errQueue := st.db.Exec(sqlQueueTable)
	if errQueue != nil {
		log.Println(errQueue)
		return errQueue
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debugEnabled bool) StoreInterface {
	st.debugEnabled = debugEnabled
	return st
}

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

func (store *Store) QueuedTaskProcess(queuedTask QueueInterface) (bool, error) {
	// 1. Start queued task
	attempts := queuedTask.Attempts() + 1

	queuedTask.AppendDetails("Task started")
	queuedTask.SetStatus(QueueStatusRunning)
	queuedTask.SetAttempts(attempts)
	queuedTask.SetStartedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	err := store.QueueUpdate(queuedTask)

	if err != nil {
		return false, err
	}

	// 2. Find task definition
	task, err := store.TaskFindByID(queuedTask.TaskID())

	if err != nil {
		return false, err
	}

	if task == nil {
		queuedTask.AppendDetails("Task DOES NOT exist")
		queuedTask.SetStatus(QueueStatusFailed)
		queuedTask.SetCompletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
		err = store.QueueUpdate(queuedTask)

		if err != nil {
			if store.debugEnabled {
				log.Println(err)
			}

			return false, err
		}

		return false, nil
	}

	handlerFunc := store.taskHandlerFunc(task.Alias())

	result := handlerFunc(queuedTask)

	if result {
		queuedTask.AppendDetails("Task completed")
		err = store.QueueSuccess(queuedTask)

		if err != nil {
			if store.debugEnabled {
				log.Println(err)
			}
		}
	} else {
		queuedTask.AppendDetails("Task failed")
		err = store.QueueFail(queuedTask)

		if err != nil {
			if store.debugEnabled {
				log.Println(err)
			}
		}
	}

	return true, nil
}

// TaskExecuteCli - CLI tool to find a task by its alias and execute its handler
// - alias "list" is reserved. it lists all the available commands
func (store *Store) TaskExecuteCli(alias string, args []string) bool {
	argumentsMap := utils.ArgsToMap(args)
	cfmt.Infoln("Executing task: ", alias, " with arguments: ", argumentsMap)

	// Lists the available tasks
	if alias == "list" {
		for index, taskHandler := range store.TaskHandlerList() {
			cfmt.Warningln(utils.ToString(index+1) + ". Task Alias: " + taskHandler.Alias())
			cfmt.Infoln("    - Task Title: " + taskHandler.Title())
			cfmt.Infoln("    - Task Description: " + taskHandler.Description())
		}

		return true
	}

	// Finds the task and executes its handler
	for _, taskHandler := range store.TaskHandlerList() {
		if strings.EqualFold(unifyName(taskHandler.Alias()), unifyName(alias)) {
			taskHandler.SetOptions(argumentsMap)
			taskHandler.Handle()
			return true
		}
	}

	cfmt.Errorln("Unrecognized task alias: ", alias)
	return false
}

func unifyName(name string) string {
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return name
}

// taskHandlerFunc finds the TaskHandler for the queued task and returns
// the Handle function, if not found, a default Handle function is returned
// which will print "No handler for alias ALIASNAME" message to notify the
// queue admin
func (store *Store) taskHandlerFunc(taskAlias string) func(queuedTask QueueInterface) bool {
	unifyName := func(name string) string {
		name = strings.ReplaceAll(name, "-", "")
		name = strings.ReplaceAll(name, "_", "")
		return name
	}

	for _, taskHandler := range store.taskHandlers {
		if strings.EqualFold(unifyName(taskHandler.Alias()), unifyName(taskAlias)) {
			return func(queuedTask QueueInterface) bool {
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

	return func(queuedTask QueueInterface) bool {
		queuedTask.AppendDetails("No handler for alias: " + taskAlias)
		store.QueueUpdate(queuedTask)
		return false
	}
}
