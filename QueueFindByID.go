package taskstore

import (
	"context"
	"database/sql"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
)

// QueueFindByID finds a Queue by ID
func (st *Store) QueueFindByID(id string) (*Queue, error) {
	sqlStr, _, err := goqu.Dialect(st.dbDriverName).
		From(st.queueTableName).
		Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNull()).
		Select().
		Limit(1).
		ToSQL()

	if err != nil {
		return nil, err
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	var Queue Queue
	err = sqlscan.Get(context.Background(), st.db, &Queue, sqlStr)

	if err != nil {
		if err == sql.ErrNoRows {
			// sqlscan does not use this anymore
			return nil, nil
		}

		if sqlscan.NotFound(err) {
			return nil, nil
		}

		log.Println("QueueStore. QueueFindByID. Error: ", err)
		return nil, err
	}

	return &Queue, nil
}
