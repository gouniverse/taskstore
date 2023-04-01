package taskstore

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func InitStore() *Store {
	db := InitDB("test_taskstore.db")
	return &Store{
		taskDefinitionTableName: "test_taskTableName.db",
		taskTaskTableName:       "Task",
		db:                      db,
		dbDriverName:            "sqlite",
		automigrateEnabled:      false,
		debug:                   false,
	}
}

func TestWithDb(t *testing.T) {
	db := InitDB("test")
	s := InitStore()

	f := WithDb(db)
	f(s)

	if s.db == nil {
		t.Fatalf("DB: Expected Initialized DB, received [%v]", s.db)
	}

}

func TestWithDefinitionTableName(t *testing.T) {
	s := InitStore()

	table_name := "test_taskTableName.db"
	f := WithDefinitionTableName(table_name)
	f(s)
	if s.taskDefinitionTableName != table_name {
		t.Fatalf("Expected DefinitionTableName [%v], received [%v]", table_name, s.taskDefinitionTableName)
	}
	table_name = "Table2"
	f = WithDefinitionTableName(table_name)
	f(s)
	if s.taskDefinitionTableName != table_name {
		t.Fatalf("Expected DefinitionTableName [%v], received [%v]", table_name, s.taskDefinitionTableName)
	}
}

func TestWithTaskTableName(t *testing.T) {
	s := InitStore()

	table_name := "test_taskTableName.db"
	f := WithTaskTableName(table_name)
	f(s)
	if s.taskTaskTableName != table_name {
		t.Fatalf("Expected TaskTableName [%v], received [%v]", table_name, s.taskTaskTableName)
	}
	table_name = "Table2"
	f = WithTaskTableName(table_name)
	f(s)
	if s.taskTaskTableName != table_name {
		t.Fatalf("Expected TaskTableName [%v], received [%v]", table_name, s.taskTaskTableName)
	}
}

func TestWithDebug(t *testing.T) {
	s := InitStore()

	b := false
	f := WithDebug(b)
	f(s)
	if s.debug != b {
		t.Fatalf("Expected Debug [%v], received [%v]", b, s.debug)
	}
}

func Test_Store_DriverName(t *testing.T) {
	db := InitDB("sqlite")
	store := InitStore()
	s := store.DriverName(db)
	if s != "sqlite" {
		t.Fatalf("Expected Debug [%v], received [%v]", "sqlite", s)
	}
}
