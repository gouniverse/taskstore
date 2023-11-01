package taskstore

// SqlCreateQueueTable returns a SQL string for creating the Queue table
func (st *Store) SqlCreateQueueTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.queueTableName + ` (
	  id             varchar(40) NOT NULL PRIMARY KEY,
	  status         varchar(40) NOT NULL,
	  task_id  varchar(40) NOT NULL,
	  parameters     text        NOT NULL,
	  output         longtext    DEFAULT NULL,
	  details        longtext    DEFAULT NULL,
	  attempts       int         DEFAULT 0,
	  started_at     datetime    DEFAULT NULL,
	  completed_at   datetime    DEFAULT NULL,
	  created_at	 datetime,
	  updated_at	 datetime,	
	  deleted_at	 datetime    DEFAULT NULL
	);
	`

	sqlPostgres := `
	CREATE TABLE IF NOT EXISTS "` + st.queueTableName + `" (
	  "id"            varchar(40)    NOT NULL PRIMARY KEY,
	  "status"        varchar(40)    NOT NULL,
	  "task_id" varchar(40)    NOT NULL,
	  "parameters"    varchar(40)    NOT NULL,
	  "output"        longtext       DEFAULT NULL,
	  "details"       longtext       DEFAULT NULL,
	  "attempts"      int            DEFAULT 0,
	  "started_at"    timestamptz(6) DEFAULT NULL,
	  "completed_at"  timestamptz(6) DEFAULT NULL,
	  "created_at"    timestamptz(6) NOT NULL,
	  "updated_at"    timestamptz(6) NOT NULL,
	  "deleted_at"    timestamptz(6) DEFAULT NULL
	)
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.queueTableName + `" (
	  "id"            varchar(40) NOT NULL PRIMARY KEY,
	  "status"        varchar(40) NOT NULL,
	  "task_id" varchar(40) NOT NULL,
	  "parameters"    varchar(40) NOT NULL,
	  "output"        text  DEFAULT NULL,
	  "details"       text  DEFAULT NULL,
	  "attempts"      int   DEFAULT 0,
	  "started_at"    datetime  DEFAULT NULL,
	  "completed_at"  datetime  DEFAULT NULL,
	  "created_at"    datetime NOT NULL,
	  "updated_at"    datetime NOT NULL,
	  "deleted_at"    datetime DEFAULT NULL
	)
	`

	sql := "unsupported driver " + st.dbDriverName

	if st.dbDriverName == "mysql" {
		sql = sqlMysql
	}
	if st.dbDriverName == "postgres" {
		sql = sqlPostgres
	}
	if st.dbDriverName == "sqlite" {
		sql = sqlSqlite
	}

	return sql
}
