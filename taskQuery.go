package taskstore

import (
	"github.com/doug-martin/goqu/v9"
)

func (st *Store) taskQuery(options TaskQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(st.dbDriverName).From(st.taskTableName)

	if options.Alias != "" {
		q = q.Where(goqu.C(COLUMN_ALIAS).Eq(options.Alias))
	}

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_TASK_ID).Eq(options.ID))
	}

	if len(options.IDIn) > 0 {
		q = q.Where(goqu.C(COLUMN_TASK_ID).In(options.IDIn))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	}

	if options.CreatedAtGreaterThan != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Gt(options.CreatedAtGreaterThan))
	}

	if options.UpdatedAtGreaterThan != "" {
		q = q.Where(goqu.C(COLUMN_UPDATED_AT).Gt(options.UpdatedAtGreaterThan))
	}

	if options.CreatedAtLessThan != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lt(options.CreatedAtLessThan))
	}

	if options.UpdatedAtLessThan != "" {
		q = q.Where(goqu.C(COLUMN_UPDATED_AT).Lt(options.UpdatedAtLessThan))
	}

	q = q.Where(goqu.C(COLUMN_DELETED_AT).IsNull())

	return q
}
