package sqltocsv

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"testing"
	"time"
)

func TestWriteCsvToFile(t *testing.T) {
	checkQueryAgainstResult(t, func(rows *sql.Rows) string {
		testCsvFileName := "/tmp/test.csv"
		err := WriteCsvToFile(testCsvFileName, rows)
		if err != nil {
			t.Fatalf("error in WriteCsvToFile: %v", err)
		}

		bytes, err := ioutil.ReadFile(testCsvFileName)
		if err != nil {
			t.Fatalf("error reading %v: %v", testCsvFileName, err)
		}

		return string(bytes[:])
	})
}

func TestWriteCsvToWriter(t *testing.T) {
	checkQueryAgainstResult(t, func(rows *sql.Rows) string {
		buffer := &bytes.Buffer{}

		err := WriteCsvToWriter(buffer, rows)
		if err != nil {
			t.Fatalf("error in WriteCsvToWriter: %v", err)
		}

		return buffer.String()
	})
}

func checkQueryAgainstResult(t *testing.T, innerTestFunc func(*sql.Rows) string) {
	db := setupDatabase(t)

	rows, err := db.Query("SELECT|people|name,age,bdate|")
	if err != nil {
		t.Fatalf("error querying: %v", err)
	}

	expectedResult := "name,age,bdate\nAlice,1,1973-11-30 08:33:09 +1100 EST\n"

	actualResult := innerTestFunc(rows)

	if actualResult != expectedResult {
		t.Errorf("Expected CSV:\n\n%v\n Got CSV:\n\n%v\n", expectedResult, actualResult)
	}
}

func setupDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("test", "foo")
	if err != nil {
		t.Fatalf("Error opening testdb %v", err)
	}
	exec(t, db, "WIPE")
	exec(t, db, "CREATE|people|name=string,age=int32,bdate=datetime")
	exec(t, db, "INSERT|people|name=Alice,age=?,bdate=?", 1, time.Unix(123456789, 0))
	return db
}

func exec(t testing.TB, db *sql.DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		t.Fatalf("Exec of %q: %v", query, err)
	}
}
