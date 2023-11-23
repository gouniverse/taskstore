package taskstore

import (
	"strings"
	"testing"
)

func Test_QueueDeleteByID(t *testing.T) {
	store, err := InitStore("test_queue_delete_by_id.db")
	if err != nil {
		t.Fatalf("QueueList: Error[%v]", err)
	}

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueList: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueList: Table creation error: [%v]", err)
	}

	queuedTask := NewQueue()
	queuedTask.Status = QueueStatusQueued

	err = store.QueueCreate(queuedTask)

	if err != nil {
		t.Fatal("QueueList: Error in creating queued task:", err.Error())
	}

	foundQueuedTask, err := store.QueueFindByID(queuedTask.ID)

	if err != nil {
		t.Fatal("QueueDeletedByID: Error in creating queued task:", err.Error())
	}

	if foundQueuedTask == nil {
		t.Fatal("QueueDeletedByID: queued task not found:")
	}

	err = store.QueueDeleteByID(queuedTask.ID)

	if err != nil {
		t.Error("QueueDeletedByID: Error deleting queued task:", err.Error())
	}

}
