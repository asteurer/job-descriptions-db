package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Job struct {
	Company     string `json:"company"`
	Position    string `json:"position"`
	Description string `json:"description"`
}

func main() {
	fileBytes, err := os.ReadFile("./data.json")
	if err != nil {
		panic(err)
	}

	var jobs []Job
	if err := json.Unmarshal(fileBytes, &jobs); err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	stmt, err := tx.Prepare("INSERT INTO jobs (company_name, job_title, job_description) VALUES(?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for _, job := range jobs {
		if _, err := stmt.Exec(job.Company, job.Position, job.Description); err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	fmt.Println("SUCCESS: File data has been transferred to the database")
}
