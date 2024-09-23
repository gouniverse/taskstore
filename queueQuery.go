package taskstore

import "github.com/doug-martin/goqu/v9"

func (st *Store) queueQuery(options QueueQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(st.dbDriverName).From(st.queueTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if options.TaskID != "" {
		q = q.Where(goqu.C(COLUMN_TASK_ID).Eq(options.TaskID))
	}

	if options.CreatedAtGte != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lt(options.CreatedAtGte))
	}

	if options.UpdatedAtLte != "" {
		q = q.Where(goqu.C(COLUMN_UPDATED_AT).Lt(options.UpdatedAtLte))
	}

	if options.SortBy != "" {
		if options.SortOrder == ASC {
			q = q.Order(goqu.I(options.SortBy).Asc())
		} else {
			q = q.Order(goqu.I(options.SortBy).Desc())
		}
	}

	q = q.Where(goqu.C(COLUMN_DELETED_AT).IsNull())

	if options.Limit != 0 {
		q = q.Limit(uint(options.Limit))
	}

	if options.Offset != 0 && !options.CountOnly {
		q = q.Offset(uint(options.Offset))
	}

	return q
}
