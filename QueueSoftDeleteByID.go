package taskstore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

func (st *Store) QueueSoftDeleteByID(id string) error {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).
		Update(st.queueTableName).
		Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNull()).
		Set(goqu.Record{"deleted_at": time.Now()}).
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
