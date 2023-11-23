package taskstore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

// QueueCreate creates a queued task
func (st *Store) QueueCreate(queue *Queue) error {
	if queue.ID == "" {
		time.Sleep(1 * time.Millisecond) // !!! important
		queue.ID = uid.MicroUid()
	}
	queue.CreatedAt = time.Now()
	queue.UpdatedAt = time.Now()

	var sqlStr string
	sqlStr, _, _ = goqu.Dialect(st.dbDriverName).Insert(st.queueTableName).Rows(queue).ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return err
	}

	return nil
}
