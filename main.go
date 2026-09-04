package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const connStr = "postgres://localhost:5432/heblang?sslmode=disable"

// setup drops and recreates the words table, then inserts one row.
func setup(db *sql.DB) error {
	if _, err := db.Exec(`DROP TABLE IF EXISTS words`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE words (id SERIAL PRIMARY KEY, hebrew TEXT)`); err != nil {
		return err
	}
	if _, err := db.Exec(`INSERT INTO words (hebrew) VALUES ($1)`, "שלום"); err != nil {
		return err
	}
	return nil
}

func main() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := setup(db); err != nil {
		log.Fatal(err)
	}
	fmt.Println("done: inserted 1 row into words")
}
