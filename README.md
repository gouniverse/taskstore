# Task Store <a href="https://gitpod.io/#https://github.com/gouniverse/taskstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>


[![Tests Status](https://github.com/gouniverse/taskstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/gouniverse/taskstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/taskstore)](https://goreportcard.com/report/github.com/gouniverse/taskstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/taskstore)](https://pkg.go.dev/github.com/gouniverse/taskstore)

TaskStore is a package to queue tasks and perform work asynchronously in the background, outside of the regular application flow.

The queue is stored in the database - SQLite, MySQL or PostgreSQL

## Installation

```
go get github.com/gouniverse/taskstore
```

## Setup

```golang
myTaskStore = taskstore.NewStore(taskstore.NewStoreOptions{
	DB:                 databaseInstance,
    	TaskTableName:      "my_task"
	QueueTableName:     "my_queue",
	AutomigrateEnabled: true,
	DebugEnabled:       false,
})
```

## Tasks

The task specifies a unit of work to be completed. It can be completed on the fly, 
or queued to the database to be completed in the background outside of
the regular application flow.

Each task has an alias (human readable identificator) that allows you to call the task,
and a title and description to give you more information on the task.

The task must implement the TaskHandlerInterface, and also define a handle method, 
which will be called to complete the queued task. 

The task may (optional) extend the TaskHandlerBase struct for additional functionality
i.e. getting task parameters, etc.

The tasks can be run directly on the terminal (CLI), or as part of a background queue.

The tasks placed in the queue will be processed at the specified interval.

```golang
package tasks

func NewHelloWorldTask() *HelloWorldTask {
	return &HelloWorldTask{}
}

type HelloWorldTask struct {
	taskstore.TaskHandlerBase
}

var _ taskstore.TaskHandlerInterface = (*HelloWorldTask)(nil) // verify it extends the task handler interface

func (task *HelloWorldTask) Alias() string {
	return "HelloWorldTask"
}

func (task *HelloWorldTask) Title() string {
	return "Hello World"
}

func (task *HelloWorldTask) Description() string {
	return "Say hello world"
}

// Enqueue. Optional shortcut to quickly add this task to the queue
func (task *HelloWorldTask) Enqueue(name string) (task *taskstore.Queue, err error) {
	return myTaskStore.TaskEnqueueByAlias(task.Alias(), map[string]any{
		"name": name,
	})
}

func (task *HelloWorldTask) Handle() bool {
	name := handler.GetParam("name")

        // Optional to allow adding the task to the queue manually. Useful while in development
	if !task.HasQueuedTask() && task.GetParam("enqueue") == "yes" {
		_, err := handler.Enqueue(name)

		if err != nil {
			task.LogError("Error enqueuing task: " + err.Error())
		} else {
			task.LogSuccess("Task enqueued.")
		}
		
		return true
	}

        if name != "" {
		task.LogInfo("Hello" + name + "!")	
	} else {
		task.LogInfo("Hello World!")
	}

	return true
}
```
## Registering the Tasks to the TaskStore

Registering the task to the task store will persist it in the database.

```
myTaskStore.TaskHandlerAdd(tasks.HelloWorldTask(), true)
```


## Adding the Task to the Queue

To add a task to the queue

```
myTaskStore.TaskEnqueueByAlias(NewHelloWorldTask.Alias(), map[string]any{
	"name": name,
})
```

Or if you have defined an Enqueue method as in the example task above.
```
NewHelloWorldTask().Enqueue("Tom Jones")
```

## Store Methods

- AutoMigrate() error - automigrate (creates) the task and queue table
- EnableDebug(debug bool) - enables / disables the debug option

## Task Methods
- TaskCreate(Task *Task) (bool, error) -  creates a Task
- TaskEnqueueByAlias(taskAlias string, parameters map[string]interface{}) (*Queue, error) -  finds a task by its alias and appends it to the queue
- TaskList(options map[string]string) ([]Task, error) - lists tasks
- TaskFindByAlias(alias string) *Task - finds a Task by alias
- TaskFindByID(id string) *Task - finds a task by ID
- TaskUpdate(Task *Task) bool - updates a task

## Queue Methods
- QueueCreate(queue *Queue) error - creates a new queued task
- QueueDeleteByID(id string) *Queue - deleted a queued task by ID
- QueueFindByID(id string) *Queue - finds a queued task by ID
- QueueFail(queue *Queue) error - fails a queued task
- QueueSoftDeleteByID(id string) *Queue - soft delete a queued task by ID (populates the deleted_at field)
- QueueSuccess(queue *Queue) error -  completes a queued task  successfully
- QueueList(options QueueListOptions) ([]Queue, error) - lists the queued tasks
- QueueUpdate(queue *Queue) error - updates a queued task
- Queue > GetParameters() (map[string]interface{}, error) - gets the parameters of the queued task
- Queue > AppendDetails(details string) - appends details to the queued task

## Similar

- https://github.com/bamzi/jobrunner
- https://github.com/rk/go-cron
- https://github.com/fieldryand/goflow
- https://github.com/go-co-op/gocron
- https://github.com/exograd/eventline
- https://github.com/ajvb/kala
