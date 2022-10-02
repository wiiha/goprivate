package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wiiha/goprivate/snote"
)

func main() {

	repo := snote.NewSqliteDefaultRepo()

	err := repo.Migrate()

	if err != nil {
		log.Fatalf("migrating posts: %v", err)
	}

}
