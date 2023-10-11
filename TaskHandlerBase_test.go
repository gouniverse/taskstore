package taskstore

import (
	"testing"
)

func Test_TaskHandlerBase_SqlCreateQueueTable(t *testing.T) {
	handler := newTestTaskHandler()

	handler.SetOptions(map[string]string{
		"completeWithSucceed": "yes",
	})

	isOK := handler.Handle()

	if !isOK {
		t.Fatalf("Test 1: Error in Handling Task")
	}

	handler2 := newTestTaskHandler()

	handler2.SetOptions(map[string]string{
		"completeWithFail": "yes",
	})

	isOK = handler2.Handle()

	if !isOK {
		t.Fatalf("Test 2: Error in Handling Task")
	}

	// store, err := InitStore("test_queue_table_create.db")
	// if err != nil {
	// 	t.Fatalf("SqlCreateQueueTable: Error[%v]", err)
	// }

	// query := store.SqlCreateQueueTable()
	// if strings.Contains(query, "unsupported driver") {
	// 	t.Fatalf("SqlCreateQueueTable: Unexpected Query, received [%v]", query)
	// }
}

func newTestTaskHandler() *testTaskHandler {
	return &testTaskHandler{}
}

type testTaskHandler struct {
	TaskHandlerBase
}

var _ TaskHandlerInterface = (*testTaskHandler)(nil) // verify it extends the task interface

func (handler *testTaskHandler) Alias() string {
	return "HelloWorldTaskHandler"
}

func (handler *testTaskHandler) Title() string {
	return "Hello World"
}

func (handler *testTaskHandler) Description() string {
	return "Say hello world"
}

// Enqueue. Optional shortcut to quickly add this task to the queue
// func (handler *testHelloWorldTaskHandler) Enqueue() (task *Queue, err error) {
// 	return myTaskStore.TaskEnqueueByAlias(handler.Alias(), map[string]any{})
// }

func (handler *testTaskHandler) Handle() bool {
	if handler.GetParam("completeWithSucceed") == "yes" {
		handler.LogSuccess("Task completed successfully.")
		return true
	}

	if handler.GetParam("completeWithFail") == "yes" {
		handler.LogError("Task completed with error.")
		return true
	}

	// 	// Optional to allow adding the task to the queue manually. Useful while in development
	// 	if !handler.HasQueuedTask() && handler.GetParam("enqueue") == "yes" {
	// 		_, err := handler.Enqueue()
	// 		if err != nil {
	// 			handler.LogError("Error enqueing task: " + err.Error())
	// 		} else {
	// 			handler.LogSuccess("Task enqueued.")
	// 		}
	// 		return true
	// 	}

	handler.LogInfo("Hello World!")
	return true
}
