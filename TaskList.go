package taskstore

import (
	"context"
	"database/sql"
	"log"

	"github.com/georgysavva/scany/sqlscan"
)

func (st *Store) TaskList(options TaskQueryOptions) ([]Task, error) {
	q := st.taskQuery(options)

	sqlStr, _, _ := q.Select().ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	list := []Task{}
	err := sqlscan.Select(context.Background(), st.db, &list, sqlStr)

	if err != nil {
		if err == sql.ErrNoRows {
			// sqlscan does not use this anymore
			return nil, nil
		}

		if sqlscan.NotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return list, err
}
