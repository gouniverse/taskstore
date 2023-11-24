package taskstore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

// TaskCreate creates a Task
func (st *Store) TaskCreate(Task *Task) (bool, error) {
	if Task.ID == "" {
		time.Sleep(1 * time.Millisecond) // !!! important
		Task.ID = uid.MicroUid()
	}
	Task.CreatedAt = time.Now()
	Task.UpdatedAt = time.Now()

	var sqlStr string
	sqlStr, _, _ = goqu.Dialect(st.dbDriverName).
		Insert(st.taskTableName).
		Rows(Task).
		ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return false, err
	}

	return true, nil
}
