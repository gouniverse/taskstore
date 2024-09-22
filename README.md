# Task Store <a href="https://gitpod.io/#https://github.com/gouniverse/taskstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>


[![Tests Status](https://github.com/gouniverse/taskstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/gouniverse/taskstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/taskstore)](https://goreportcard.com/report/github.com/gouniverse/taskstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/taskstore)](https://pkg.go.dev/github.com/gouniverse/taskstore)

TaskStore is a robust, asynchronous durable task queue package designed to offload time-consuming or resource-intensive operations from your main application.

By deferring tasks to the background, you can improve application responsiveness and prevent performance bottlenecks.

TaskStore leverages a durable database (SQLite, MySQL, or PostgreSQL) to ensure reliable persistence and fault tolerance.

## License

This project is licensed under the GNU General Public License version 3 (GPL-3.0). You can find a copy of the license at https://www.gnu.org/licenses/gpl-3.0.en.html

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

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

The task specifies a unit of work to be completed. It can be performed immediately, 
or enqueued to the database and deferreed for asynchronious processing, ensuring your
application remains responsive.

Each task is uniquely identified by an alias and provides a human-readable title and description.

Each task is uniquely identified by an alias that allows the task to be easily called. 
A human-readable title and description to give the user more information on the task.

To define a task, implement the TaskHandlerInterface and provide a Handle method
that contains the task's logic.

Optionally, extend the TaskHandlerBase struct for additional features like parameter
retrieval.

Tasks can be executed directly from the command line (CLI) or as part of a background queue.

The tasks placed in the queue will be processed at specified interval.

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
myTaskStore.TaskHandlerAdd(NewHelloWorldTask(), true)
```

## Executing Tasks in the Terminal

To add the option to execute tasks from the terminal add the following to your main method

```
myTaskStore.TaskExecuteCli(args[1], args[1:])
```

Example:
```
go run . HelloWorldTask --name="Tom Jones"
```

## Adding the Task to the Queue

To add a task to the background queue

```
myTaskStore.TaskEnqueueByAlias(NewHelloWorldTask.Alias(), map[string]any{
	"name": name,
})
```

Or if you have defined an Enqueue method as in the example task above.
```
NewHelloWorldTask().Enqueue("Tom Jones")
```

## Starting the Queue

To start the queue call the QueueRunGoroutine. 
It allows you to specify the interval for processing the queued tasks (i.e. every 10 seconds)
Also to set timeout for queued tasks. After a queued task is started if it has not completed in the specified timeout it will be marked as failed, and the rest of he tasks will start to be processed.

```
myTaskStore.QueueRunGoroutine(10, 2) // every 10s, unstuck after 2 mins
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

## Frequently Asked Questions (FAQ)

### 1. What is TaskStore used for?
TaskStore is a versatile tool for offloading time-consuming or resource-intensive
tasks from your main application. By deferring these tasks to the background,
you can improve application responsiveness and prevent performance bottlenecks.

It's ideal for tasks like data processing, sending emails, generating reports,
or performing batch operations.

### 2. How does TaskStore work?
TaskStore creates a durable queue in your database (SQLite, MySQL, or PostgreSQL)
to store tasks. These tasks are then processed asynchronously by a background worker.
You can define tasks using a simple interface and schedule them to be executed
at specific intervals or on demand.

### 3. What are the benefits of using TaskStore?

- Improved application performance: Offload time-consuming tasks to prevent performance bottlenecks.
- Asynchronous processing: Execute tasks independently of your main application flow.
- Reliability: Ensure tasks are completed even if your application crashes.
- Flexibility: Schedule tasks to run at specific intervals or on demand.
- Ease of use: Define tasks using a simple interface and integrate with your existing application.

### 4. How do I create a task in TaskStore?
To create a task, you'll need to implement the TaskHandlerInterface and provide a Handle method that contains the task's logic. You can also extend the TaskHandlerBase struct for additional features.

### 5. How do I schedule a task to run in the background?
Use the TaskEnqueueByAlias method to add a task to the background queue. You can specify the interval at which the queue is processed using the QueueRunGoroutine method.

### 6. Can I monitor the status of tasks?
Yes, TaskStore provides methods to list tasks, check their status, and view task details.

### 7. How does TaskStore handle task failures?
If a task fails, it can be retried automatically or marked as failed. You can customize the retry logic to suit your specific needs.

### 8. Is TaskStore suitable for large-scale applications?
Yes, TaskStore is designed to handle large volumes of tasks. It can be scaled horizontally by adding more worker processes.

### 9. Does TaskStore support different database systems?
Yes, TaskStore supports SQLite, MySQL, and PostgreSQL.

### 10. Can I customize TaskStore to fit my specific needs?
Yes, TaskStore is highly customizable. You can extend and modify the code to suit your requirements.

## Similar

- https://github.com/harshadmanglani/polaris
- https://github.com/bamzi/jobrunner
- https://github.com/rk/go-cron
- https://github.com/fieldryand/goflow
- https://github.com/go-co-op/gocron
- https://github.com/exograd/eventline
- https://github.com/ajvb/kala
- https://github.com/shiblon/taskstore
