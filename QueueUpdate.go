package taskstore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// QueueUpdate creates a Queue
func (st *Store) QueueUpdate(queue *Queue) error {
	queue.UpdatedAt = time.Now()

	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).
		Update(st.queueTableName).
		Where(goqu.C(COLUMN_ID).Eq(queue.ID)).
		Set(queue).
		ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		if st.debugEnabled {
			log.Println(err)
		}

		return err
	}

	return nil
}
