package taskstore

import "github.com/gouniverse/sb"

// SqlCreateTaskTable returns a SQL string for creating the Task table
func (st *Store) SqlCreateTaskTable() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(st.db)).
		Table(st.taskTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     40,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_ALIAS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 100,
		}).
		Column(sb.Column{
			Name:   COLUMN_TITLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: COLUMN_DESCRIPTION,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     COLUMN_DELETED_AT,
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	return sql
}
