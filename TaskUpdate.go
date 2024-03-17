package taskstore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// TaskUpdate updates a Task
func (st *Store) TaskUpdate(task *Task) error {
	task.UpdatedAt = time.Now()

	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).
		Update(st.queueTableName).
		Where(goqu.C(COLUMN_ID).Eq(task.ID)).
		Set(task).
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
