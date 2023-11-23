package taskstore

import (
	"context"
	"database/sql"
	"log"

	"github.com/georgysavva/scany/sqlscan"
)

func (st *Store) QueueList(options QueueQueryOptions) ([]Queue, error) {
	q := st.queueQuery(options)

	sqlStr, _, _ := q.Select().ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	list := []Queue{}
	err := sqlscan.Select(context.Background(), st.db, &list, sqlStr)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		}

		if sqlscan.NotFound(err) {
			return nil, nil
		}

		log.Println("QueueStore. QueueList. Error: ", err)
		return nil, err
	}

	return list, err
}
