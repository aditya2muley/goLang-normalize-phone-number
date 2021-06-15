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

func checkError(err error) error {
	if err != nil {
		panic(err)
	}
	return nil
}
func main() {
	dbInfo := fmt.Sprintf("user=%s password=%d host=%s port=%d sslmode=disable", username, password, host, port)
	db, err := sql.Open("postgres", dbInfo)
	checkError(err)
	fmt.Println(dbInfo)
	err = resetDB(db, dbName)
	checkError(err)
	db.Close()
	dbInfo = fmt.Sprintf("%s dbname=%s", dbInfo, dbName)
	db, err = sql.Open("postgres", dbInfo)
	checkError(err)
	err = createTable(db)
	checkError(err)
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
	checkError(err)
	return nil
}

func resetDB(db *sql.DB, dbName string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + dbName)
	checkError(err)
	return createDB(db, dbName)
}
