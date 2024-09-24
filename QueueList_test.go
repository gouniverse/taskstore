package taskstore

import (
	"strings"
	"testing"
)

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

	list, err := store.QueueList(QueueQueryOptions{
		Status:    QueueStatusQueued,
		Limit:     10,
		SortOrder: ASC,
		SortBy:    COLUMN_CREATED_AT,
	})

	if err != nil {
		t.Fatalf("QueueList: Error[%v]", err)
	}

	if len(list) != 1 {
		t.Fatal("There must be 1 task, found: ", list)
	}
}
