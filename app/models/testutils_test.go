package models

import (
	"database/sql"
	"io/ioutil"
	"testing"

	_ "github.com/lib/pq"
)

const dbScriptPath = "../../database/"

// newTestDB returns the mock database connection along with its teardown function.
func newTestDB(t *testing.T) (*sql.DB, func()) {
	// Database connection.
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=password dbname=memebot_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	// Database setup.
	script, err := ioutil.ReadFile(dbScriptPath + "setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Insert mock data.
	script, err = ioutil.ReadFile(dbScriptPath + "mock_data.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Return the db and the tear down function.
	return db, func() {
		script, err := ioutil.ReadFile(dbScriptPath + "teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	}
}
