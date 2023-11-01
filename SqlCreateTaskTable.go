package taskstore

// SqlCreateTaskTable returns a SQL string for creating the Task table
func (st *Store) SqlCreateTaskTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.taskTableName + ` (
	  id             varchar(40)  NOT NULL PRIMARY KEY,
	  status         varchar(40)  NOT NULL,
	  alias          varchar(100) NOT NULL,
	  title          varchar(255) NOT NULL,
	  description    text         DEFAULT NULL,
	  memo           text         DEFAULT NULL,
	  created_at	 datetime     NOT NULL,
	  updated_at	 datetime 	  NOT NULL,
	  deleted_at	 datetime     DEFAULT NULL
	);
	`

	sqlPostgres := `
	CREATE TABLE IF NOT EXISTS "` + st.taskTableName + `" (
	  "id"          varchar(40)    NOT NULL PRIMARY KEY,
	  "status"      varchar(40)    NOT NULL,
	  "alias"       varchar(40)    NOT NULL,
	  "title"       varchar(255)   NOT NULL,
	  "description" text,          DEFAULT NULL
	  "memo"        text,          DEFAULT NULL
	  "created_at"  timestamptz(6) NOT NULL,
	  "updated_at"  timestamptz(6) NOT NULL,
	  "deleted_at"  timestamptz(6) DEFAULT NULL
	);
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.taskTableName + `" (
	  "id"          varchar(40)  NOT NULL PRIMARY KEY,
	  "status"      varchar(40)  NOT NULL,
	  "alias"       varchar(40)  NOT NULL,
	  "title"       varchar(255) NOT NULL,
	  "description" text,
	  "memo"        text,
	  "created_at"  datetime     NOT NULL,
	  "updated_at"  datetime     NOT NULL,
	  "deleted_at"  datetime     DEFAULT NULL
	);
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
