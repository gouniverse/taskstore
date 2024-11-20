package taskstore

import (
	"testing"

	"github.com/mingrammer/cfmt"
)

func Test_TaskHandlerBase(t *testing.T) {
	handler := newTestTaskHandler()

	handler.SetOptions(map[string]string{
		"completeWithSuccess": "yes",
	})

	isOK := handler.Handle()

	if !isOK {
		t.Fatalf("Test 1: Must be successful")
	}

	if handler.SuccessMessage() != "Task forced to succeed." {
		t.Fatalf("Test 1: Message must be 'Task forced to succeed.', but found: %s", handler.SuccessMessage())
	}

	handler2 := newTestTaskHandler()

	handler2.SetOptions(map[string]string{
		"completeWithFail": "yes",
	})

	isOK = handler2.Handle()

	if isOK {
		t.Fatalf("Test 2: Must Fail")
	}

	if handler2.ErrorMessage() != "Task forced to fail." {
		t.Fatalf("Test 2: Message must be 'Task forced to fail.', but found: %s", handler2.ErrorMessage())
	}
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

func (handler *testTaskHandler) Handle() bool {
	cfmt.Warningln("Param 1", handler.GetParam("completeWithSuccess"))
	cfmt.Warningln("Param 2", handler.GetParam("completeWithFail"))

	if handler.GetParam("completeWithSuccess") == "yes" {
		handler.LogSuccess("Task forced to succeed.")
		return true
	}

	if handler.GetParam("completeWithFail") == "yes" {
		handler.LogError("Task forced to fail.")
		return false
	}

	return false
}
