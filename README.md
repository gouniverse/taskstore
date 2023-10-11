# Task Store <a href="https://gitpod.io/#https://github.com/gouniverse/taskstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>


[![Tests Status](https://github.com/gouniverse/taskstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/gouniverse/taskstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/taskstore)](https://goreportcard.com/report/github.com/gouniverse/taskstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/taskstore)](https://pkg.go.dev/github.com/gouniverse/taskstore)

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

## Task Handlers

Task handlers process the queued tasks. They must implement the TaskHandlerInterface, 
and optionally extend the TaskHandlerBase struct for additional functionality
i.e. getting task parameters, etc.

The task handlers can be run directly on the CLI, or as part of a background queue.

```golang
package taskhandlers

func NewHelloWorldTaskHandler() *HelloWorldTaskHandler {
	return &HelloWorldTaskHandler{}
}

type HelloWorldTaskHandler struct {
	taskstore.TaskHandlerBase
}

var _ taskstore.TaskHandlerInterface = (*HelloWorldTaskHandler)(nil) // verify it extends the task interface

func (handler *HelloWorldTaskHandler) Alias() string {
	return "HelloWorldTaskHandler"
}

func (handler *HelloWorldTaskHandler) Title() string {
	return "Hello World"
}

func (handler *HelloWorldTaskHandler) Description() string {
	return "Say hello world"
}

// Enqueue. Optional shortcut to quickly add this task to the queue
func (handler *HelloWorldTaskHandler) Enqueue() (task *taskstore.Queue, err error) {
	return myTaskStore.TaskEnqueueByAlias(handler.Alias(), map[string]any{})
}

func (handler *HelloWorldTaskHandler) Handle() bool {

        // Optional to allow adding the task to the queue manually. Useful while in development
	if !handler.HasQueuedTask() && handler.GetParam("enqueue") == "yes" {
		_, err := handler.Enqueue()

		if err != nil {
			handler.LogError("Error enqueuing task: " + err.Error())
		} else {
			handler.LogSuccess("Task enqueued.")
		}
		
		return true
	}

	handler.LogInfo("Hello World!")
	return true
}
```

## Methods

- AutoMigrate() error - automigrate (creates) the task and queue table
- EnableDebug(debug bool) - enables / disables the debug option
- TaskCreate(Task *Task) (bool, error) -  creates a Task
- TaskEnqueueByAlias(taskAlias string, parameters map[string]interface{}) (*Queue, error) -  finds a task by its alias and appends it to the queue
- TaskList(options map[string]string) ([]Task, error) - lists tasks
- TaskFindByAlias(alias string) *Task - finds a Task by alias
- TaskFindByID(id string) *Task - finds a task by ID
- TaskUpdate(Task *Task) bool - updates a task
- QueueFail(queue *Queue) error - fails a queued task
- QueueSuccess(queue *Queue) error -  completes a queued task  successfully
- QueueCreate(queue *Queue) error - creates a new queued task
- QueueList(options QueueListOptions) ([]Queue, error) - lists the queued tasks
- QueueFindByID(id string) *Queue - finds a queued task by ID
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
