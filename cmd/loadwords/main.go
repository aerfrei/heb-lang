// Command loadwords parses tanach.us book XML files and loads their words
// into Postgres, one row per word (book, chapter, verse, position, letters).
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"

	"github.com/aerfrei/heb-lang/internal/tanach"
)

const (
	connStr  = "postgres://localhost:5432/heblang?sslmode=disable"
	booksDir = "tanach"
)

func ensureTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tanach_words (
		book     TEXT NOT NULL,
		chapter  INT  NOT NULL,
		verse    INT  NOT NULL,
		position INT  NOT NULL,
		letters  TEXT NOT NULL
	)`)
	return err
}

func loadBook(db *sql.DB, path, book string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	words, err := tanach.ParseBook(f, book)
	if err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	if _, err := tx.Exec(`DELETE FROM tanach_words WHERE book = $1`, book); err != nil {
		tx.Rollback()
		return 0, err
	}
	stmt, err := tx.Prepare(`INSERT INTO tanach_words (book, chapter, verse, position, letters) VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()
	for _, w := range words {
		if _, err := stmt.Exec(w.Book, w.Chapter, w.Verse, w.Position, w.Letters); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	return len(words), tx.Commit()
}

func main() {
	book := flag.String("book", "", "process a single book (filename without .xml, e.g. Ruth); if empty, processes every book in tanach/")
	flag.Parse()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := ensureTable(db); err != nil {
		log.Fatal(err)
	}

	var names []string
	if *book != "" {
		names = []string{*book}
	} else {
		matches, err := filepath.Glob(filepath.Join(booksDir, "*.xml"))
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range matches {
			names = append(names, strings.TrimSuffix(filepath.Base(m), ".xml"))
		}
	}

	for _, name := range names {
		path := filepath.Join(booksDir, name+".xml")
		n, err := loadBook(db, path, name)
		if err != nil {
			log.Fatalf("%s: %v", name, err)
		}
		fmt.Printf("%s: loaded %d words\n", name, n)
	}
}
