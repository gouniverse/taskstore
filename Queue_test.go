package taskstore

import (
	"encoding/json"
	"strings"
	"testing"
)

func NewQueue() *Queue {
	return &Queue{}
}

func Test_Store_SqlCreateQueueTable(t *testing.T) {
	store, err := InitStore("test_queue_table_create.db")
	if err != nil {
		t.Fatalf("SqlCreateQueueTable: Error[%v]", err)
	}

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("SqlCreateQueueTable: Unexpected Query, received [%v]", query)
	}
}

func Test_Store_QueueCreate(t *testing.T) {
	store, err := InitStore("test_queue_create.db")
	if err != nil {
		t.Fatalf("QueueCreate: Error[%v]", err)
	}

	task := NewQueue()
	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueCreate: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueCreate: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("QueueCreate: Error in Creating Queue: received [%v]", err)
	}
}

func Test_Store_QueueUpdate(t *testing.T) {
	store, err := InitStore("test_queue_update.db")
	if err != nil {
		t.Fatalf("QueueUpdate: Error[%v]", err)
	}

	task := NewQueue()
	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueUpdate: UnExpected Query, received [%v]", query)
	}
	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueUpdate: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("QueueUpdate: Error in Creating Queue: received [%v]", err)
	}

	err = store.QueueUpdate(task)
	if err != nil {
		t.Fatalf("QueueUpdate: Error in Updating Queue: received [%v]", err)
	}
}

func Test_Store_QueueSuccess(t *testing.T) {
	store, err := InitStore("test_queue_success.db")
	if err != nil {
		t.Fatalf("QueueSuccess: Error[%v]", err)
	}

	task := NewQueue()
	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueSuccess: UnExpected Query, received [%v]", query)
	}
	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueSuccess: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("QueueSuccess: Error in Creating Queue: received [%v]", err)
	}

	err = store.QueueSuccess(task)
	if err != nil {
		t.Fatalf("QueueSuccess: Error in Success Queue: received [%v]", err)
	}
}

func Test_Store_QueueFail(t *testing.T) {
	store, err := InitStore("test_queue_fail.db")
	if err != nil {
		t.Fatalf("QueueFail: Error[%v]", err)
	}

	task := NewQueue()
	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueFail: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueFail: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("QueueFail: Error in Creating Queue: received [%v]", err)
	}

	err = store.QueueFail(task)
	if err != nil {
		t.Fatalf("QueueFail: Error in Fail Queue: received [%v]", err)
	}
}

func Test_Store_QueueFindByID(t *testing.T) {
	store, err := InitStore("test_queue_find_by_id.db")
	if err != nil {
		t.Fatalf("QueueFindByID: Error[%v]", err)
	}
	task := NewQueue()
	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueFindByID: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueFindByID: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("QueueFindByID: Error in Creating Queue: received [%v]", err)
	}

	id := task.ID
	queue := store.QueueFindByID(id)
	if queue == nil {
		t.Fatalf("QueueFindByID: Error in Finding Queue: ID [%v]", id)
	}
	if queue.ID != id {
		t.Fatalf("QueueFindByID: ID not matching, Expected[%v], Received[%v]", id, queue.ID)
	}
}

func Test_Queue_AppendDetails(t *testing.T) {
	task := NewQueue()
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

func Test_Queue_GetParameters(t *testing.T) {
	store, err := InitStore("test_queue_get_parameters.db")
	if err != nil {
		t.Fatalf("GetParameters: Error[%v]", err)
	}

	task := NewQueue()
	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("GetParameters: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("GetParameters: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("GetParameters: Error in Creating Queue: received [%v]", err)
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

func Test_Queue_QueueList(t *testing.T) {
	store, err := InitStore("test_queue_list.db")
	if err != nil {
		t.Fatalf("QueueList: Error[%v]", err)
	}

	task := NewQueue()
	task.Status = QueueStatusQueued
	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueList: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueList: Table creation error: [%v]", err)
	}
	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("QueueList: Error in Creating Queue: received [%v]", err)
	}

	// u, err := json.Marshal(Temp{Status: "Bob", Limit: 10})
	// if err != nil {
	// 	t.Fatalf("%v", err)
	// }
	// task.Parameters = string(u)

	// data, err := task.GetParameters()
	// if err != nil {
	// 	t.Fatalf("QueueList: Error[%v]", err)
	// }
	list, err := store.QueueList(QueueListOptions{
		Status:    QueueStatusQueued,
		Limit:     10,
		SortOrder: "asc",
		SortBy:    "id",
	})

	if err != nil {
		t.Fatalf("QueueList: Error[%v]", err)
	}

	if len(list) != 1 {
		t.Fatal("There must be 1 task, found: ", list)
	}
}
