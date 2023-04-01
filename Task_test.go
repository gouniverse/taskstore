package taskstore

import (
	"encoding/json"
	"strings"
	"testing"
)

func NewTask() *Task {
	return &Task{}
}

func Test_Store_SqlCreateTaskTable(t *testing.T) {
	s := InitStore()
	query := s.SqlCreateTaskTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("SqlCreateTaskTable: UnExpected Query, received [%v]", query)
	}
}

func Test_Store_TaskCreate(t *testing.T) {
	s := InitStore()
	task := NewTask()
	query := s.SqlCreateTaskTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskCreate: UnExpected Query, received [%v]", query)
	}
	_, err := s.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskCreate: Table creation error: [%v]", err)
	}
	err = s.TaskCreate(task)
	if err != nil {
		t.Fatalf("TaskCreate: Error in Creating Task: received [%v]", err)
	}
}

func Test_Store_TaskUpdate(t *testing.T) {
	s := InitStore()
	task := NewTask()
	query := s.SqlCreateTaskTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskUpdate: UnExpected Query, received [%v]", query)
	}
	_, err := s.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskUpdate: Table creation error: [%v]", err)
	}
	err = s.TaskCreate(task)
	if err != nil {
		t.Fatalf("TaskUpdate: Error in Creating Task: received [%v]", err)
	}
	err = s.TaskUpdate(task)
	if err != nil {
		t.Fatalf("TaskUpdate: Error in Updating Task: received [%v]", err)
	}
}

func Test_Store_TaskSuccess(t *testing.T) {
	s := InitStore()
	task := NewTask()
	query := s.SqlCreateTaskTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskSuccess: UnExpected Query, received [%v]", query)
	}
	_, err := s.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskSuccess: Table creation error: [%v]", err)
	}
	err = s.TaskCreate(task)
	if err != nil {
		t.Fatalf("TaskSuccess: Error in Creating Task: received [%v]", err)
	}
	err = s.TaskSuccess(task)
	if err != nil {
		t.Fatalf("TaskSuccess: Error in Success Task: received [%v]", err)
	}
}

func Test_Store_TaskFail(t *testing.T) {
	s := InitStore()
	task := NewTask()
	query := s.SqlCreateTaskTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskFail: UnExpected Query, received [%v]", query)
	}
	_, err := s.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskFail: Table creation error: [%v]", err)
	}
	err = s.TaskCreate(task)
	if err != nil {
		t.Fatalf("TaskFail: Error in Creating Task: received [%v]", err)
	}
	err = s.TaskFail(task)
	if err != nil {
		t.Fatalf("TaskFail: Error in Fail Task: received [%v]", err)
	}
}

func Test_Store_TaskFindByID(t *testing.T) {
	s := InitStore()
	task := NewTask()
	query := s.SqlCreateTaskTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskFindByID: UnExpected Query, received [%v]", query)
	}
	_, err := s.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskFindByID: Table creation error: [%v]", err)
	}
	err = s.TaskCreate(task)
	if err != nil {
		t.Fatalf("TaskFindByID: Error in Creating Task: received [%v]", err)
	}
	id := task.ID
	newtask := s.TaskFindByID(id)
	if newtask == nil {
		t.Fatalf("TaskFindByID: Error in Finding Task: ID [%v]", id)
	}
	if newtask.ID != id {
		t.Fatalf("TaskFindByID: ID not matching, Expected[%v], Received[%v]", id, newtask.ID)
	}
}

func Test_Task_AppendDetails(t *testing.T) {
	task := NewTask()
	str := "Test1"
	task.AppendDetails(str)

	if !strings.Contains(task.Details, str) {
		t.Fatalf("AppendDetails: Failed Details[%v]", task.Details)
	}
}

type Temp struct {
	Status     string `json:"status"`
	Limit      int    `json:"limit"`
	Sort_by    string `json:"sort_by"`
	Sort_order string `json:"sort_order"`
}

func Test_Task_GetParameters(t *testing.T) {
	s := InitStore()
	task := NewTask()
	query := s.SqlCreateTaskTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("GetParameters: UnExpected Query, received [%v]", query)
	}
	_, err := s.db.Exec(query)
	if err != nil {
		t.Fatalf("GetParameters: Table creation error: [%v]", err)
	}
	err = s.TaskCreate(task)
	if err != nil {
		t.Fatalf("GetParameters: Error in Creating Task: received [%v]", err)
	}

	u, err := json.Marshal(Temp{Status: "Bob", Limit: 10})
	if err != nil {
		t.Fatalf("%v", err)
	}
	task.Parameters = string(u)

	_, err = task.GetParameters()
	if err != nil {
		t.Fatalf("GetParameters: Error[%v]", err)
	}
}

func Test_Task_TaskList(t *testing.T) {
	s := InitStore()
	task := NewTask()
	task.Status = TaskStatusQueued
	query := s.SqlCreateTaskTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskList: UnExpected Query, received [%v]", query)
	}
	_, err := s.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskList: Table creation error: [%v]", err)
	}
	err = s.TaskCreate(task)
	if err != nil {
		t.Fatalf("TaskList: Error in Creating Task: received [%v]", err)
	}

	// u, err := json.Marshal(Temp{Status: "Bob", Limit: 10})
	// if err != nil {
	// 	t.Fatalf("%v", err)
	// }
	// task.Parameters = string(u)

	// data, err := task.GetParameters()
	// if err != nil {
	// 	t.Fatalf("TaskList: Error[%v]", err)
	// }
	list, err := s.TaskList(TaskListOptions{
		Status:    TaskStatusQueued,
		Limit:     10,
		SortOrder: "asc",
		SortBy:    "id",
	})

	if err != nil {
		t.Fatalf("TaskList: Error[%v]", err)
	}

	if len(list) != 1 {
		t.Fatal("There must be 1 task, found: ", list)
	}
}
