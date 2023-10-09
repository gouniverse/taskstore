# Task Store <a href="https://gitpod.io/#https://github.com/gouniverse/taskstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>


[![Tests Status](https://github.com/gouniverse/taskstore/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/gouniverse/taskstore/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/taskstore)](https://goreportcard.com/report/github.com/gouniverse/taskstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/taskstore)](https://pkg.go.dev/github.com/gouniverse/taskstore)

```
go get github.com/gouniverse/taskstore
```

## Setup

```go
taskStore = taskstore.NewStore(taskstore.NewStoreOptions{
	DB:                 databaseInstance,
    TaskTableName:      "my_task"
	QueueTableName:     "my_queue",
	AutomigrateEnabled: true,
	DebugEnabled:       false,
})
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
