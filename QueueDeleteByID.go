package taskstore

import (
	"log"

	"github.com/doug-martin/goqu/v9"
)

func (st *Store) QueueDeleteByID(id string) error {
	sqlStr, preparedArgs, _ := goqu.Dialect(st.dbDriverName).
		From(st.queueTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		Delete().
		ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr, preparedArgs...)

	if err != nil {
		if st.debugEnabled {
			log.Println(err)
		}

		return err
	}

	return nil
}
