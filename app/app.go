package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	username = "test_user"
	password = 123456
	dbName   = "normalize_phone_numbers"
)

func normalizePhoneNumber(number string) string {
	reg := regexp.MustCompile("\\D")
	str := reg.ReplaceAllString(number, "")
	return str
}

func main() {
	dbInfo := fmt.Sprintf("user=%s password=%d host=%s port=%d sslmode=disable", username, password, host, port)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(dbInfo)
	err = resetDB(db, dbName)
	if err != nil {
		panic(err)
	}
	db.Close()
	dbInfo = fmt.Sprintf("%s dbname=%s", dbInfo, dbName)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	err = createTable(db)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func createTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers(
			id SERIAL,
			number VARCHAR(255)
		)
	`
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func createDB(db *sql.DB, dbName string) error {
	_, err := db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, dbName string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		return err
	}
	return createDB(db, dbName)
}
