package sqlite

import (
	"os"
	"path/filepath"
	"testing"
)

var path string = "test_data/db.db"

func TestInitDB(t *testing.T) {

	os.MkdirAll(filepath.Dir(path), 0770)

	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	if err := db.Init(); err != nil {
		t.Fatal("can't init db:", err)
	}

}

func TestCreateMockContent(t *testing.T) {
	TestInitDB(t)
	TestAddUsers(t)
	TestAddServers(t)
	TestAddSubscriptions(t)
}

func TestDropDB(t *testing.T) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatal(err)
	}
}
