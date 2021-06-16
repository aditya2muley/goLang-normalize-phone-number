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

type phone struct {
	id     int
	number string
}

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
	// err = resetDB(db, dbName)
	// checkError(err)
	db.Close()
	dbInfo = fmt.Sprintf("%s dbname=%s", dbInfo, dbName)
	db, err = sql.Open("postgres", dbInfo)
	checkError(err)
	err = createTable(db)
	checkError(err)
	numbers := []string{"1234567890", "123 456 7891", "123 456 7892", "123 456-7893", "123-456-7894", "123-456-7890", "1234567892", "123-456-7892"}
	for _, number := range numbers {
		_, err := insertRow(db, number)
		checkError(err)
	}
	data, err := viewAllRecord(db)
	checkError(err)
	data, err = viewRecord(db, 1)
	checkError(err)
	fmt.Printf("%+v", data)
	defer db.Close()
}

func viewAllRecord(db *sql.DB) ([]phone, error) {
	var ret []phone
	statement := "SELECT * FROM phone_numbers"
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var number phone
		if err := rows.Scan(&number.id, &number.number); err != nil {
			return nil, err
		}
		ret = append(ret, number)
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func viewRecord(db *sql.DB, id int) (string, error) {
	var number string
	statement := "SELECT number FROM phone_numbers WHERE id = $1"
	err := db.QueryRow(statement, id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func insertRow(db *sql.DB, num string) (int, error) {
	var id int
	statement := "INSERT INTO phone_numbers(number) values($1) RETURNING ID"
	err := db.QueryRow(statement, num).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
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
