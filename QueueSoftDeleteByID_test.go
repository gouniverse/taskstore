package taskstore

import (
	"strings"
	"testing"
)

func Test_Store_QueueSoftDeleteByID(t *testing.T) {
	store, err := InitStore("test_queue_soft_delete_by_id.db")
	if err != nil {
		t.Fatalf("QueueSoftDeleteByID: Error[%v]", err)
	}

	queuedTask := NewQueue()
	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueSoftDeleteByID: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueSoftDeleteByID: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(queuedTask)
	if err != nil {
		t.Fatalf("QueueSoftDeleteByID: Error in Creating Queue: received [%v]", err)
	}

	err = store.QueueSoftDeleteByID(queuedTask.ID)
	if err != nil {
		t.Fatalf("QueueSoftDeleteByID: Error in Fail Queue: received [%v]", err)
	}

	queueFound, err := store.QueueFindByID(queuedTask.ID)

	if err != nil {
		t.Fatal("QueueSoftDeleteByID: Error in QueueFindByID: received:", err)
	}

	if queueFound != nil {
		t.Fatal("QueueSoftDeleteByID: QueueFindByID should be nil, received:", err)
	}
}
