package taskstore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

func (st *Store) QueueSoftDeleteByID(id string) error {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).
		Update(st.queueTableName).
		Where(goqu.C(COLUMN_ID).Eq(id), goqu.C(COLUMN_DELETED_AT).IsNull()).
		Set(goqu.Record{COLUMN_DELETED_AT: time.Now()}).
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
