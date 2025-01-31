package sqlite

import (
	"encoding/json"
	"fmt"
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

// func TestDropDB(t *testing.T) {
// 	err := os.RemoveAll(path)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

func makeMap(in interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}

func compareStructs(want interface{}, got interface{}, fields []string) error {
	wantMap := makeMap(want)
	gotMap := makeMap(got)

	for _, field := range fields {
		wantField, ok1 := wantMap[field]
		gotField, ok2 := gotMap[field]
		if ok1 == ok2 && wantField == gotField {
			continue
		} else {
			return fmt.Errorf(`field "%s": want "%v", got "%v"`, field, wantField, gotField)
		}
	}

	return nil
}
