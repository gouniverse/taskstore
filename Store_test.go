package taskstore

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
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
	db := InitDB("test_settingstore.db")
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

/*
func Test_Store_AutoMigrate(t *testing.T) {
	db := InitDB("test_settingsTableName.db")

	s, _ := NewStore(WithDb(db), WithTableName("log_with_automigrate"), WithAutoMigrate(true))

	s.AutoMigrate()

	if s.settingsTableName != "log_with_automigrate" {
		t.Fatalf("Expected logTableName [log_with_automigrate] received [%v]", s.settingsTableName)
	}
	if s.db == nil {
		t.Fatalf("DB Init Failure")
	}
	if s.automigrateEnabled != true {
		t.Fatalf("Failure:  WithAutoMigrate")
	}
}
*/
