package db

import (
	"database/sql"
	"fmt"
	"regexp"
)

type DBObject struct {
	db *sql.DB
}

type Phone struct {
	Id     int
	Number string
}

func checkError(err error) error {
	if err != nil {
		panic(err)
	}
	return nil
}

func OpenConnection(driver string, credential string) *sql.DB {
	db, err := sql.Open(driver, credential)
	checkError(err)
	return db
}

func ResetDB(db *sql.DB, dbName string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + dbName)
	checkError(err)
	return createDB(db, dbName)
}

func createDB(db *sql.DB, dbName string) error {
	_, err := db.Exec("CREATE DATABASE " + dbName)
	checkError(err)
	return db.Close()
}

func CreateTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers(
			id SERIAL,
			number VARCHAR(255)
		)
	`
	_, err := db.Exec(statement)
	checkError(err)
	return nil
}

func SeedData(db *sql.DB) {
	numbers := []string{"1234567890", "123 456 7891", "123 456 7892", "123 456-7893", "123-456-7894", "123-456-7890", "1234567892", "123-456-7892"}
	for _, number := range numbers {
		_, err := insertRow(db, number)
		checkError(err)
	}
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

func ViewAllRecord(db *sql.DB) ([]Phone, error) {
	var ret []Phone
	statement := "SELECT * FROM phone_numbers"
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var number Phone
		if err := rows.Scan(&number.Id, &number.Number); err != nil {
			return nil, err
		}
		ret = append(ret, number)
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func ViewRecord(db *sql.DB, id int) (string, error) {
	var number string
	statement := "SELECT number FROM phone_numbers WHERE id = $1"
	err := db.QueryRow(statement, id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func normalizePhoneNumber(number string) string {
	reg := regexp.MustCompile("\\D")
	str := reg.ReplaceAllString(number, "")
	return str
}

func NormalizeAndUpdateRecord(db *sql.DB, number_list []Phone) {
	for _, ph := range number_list {
		number := normalizePhoneNumber(ph.Number)
		if number != ph.Number {
			fmt.Println("Changes required for phone number=", ph.Number)
			existing, err := getPhoneRecord(db, number)
			checkError(err)
			if existing != nil {
				checkError(deleteRecord(db, existing))
				ph.Number = number
				checkError(updateRecord(db, &ph))
			} else {
				ph.Number = number
				checkError(updateRecord(db, &ph))
			}
		} else {
			fmt.Println("No changes required")
		}
	}
}
func getPhoneRecord(db *sql.DB, number string) (*Phone, error) {
	var ph Phone
	statement := `SELECT * FROM phone_numbers where number = $1`
	err := db.QueryRow(statement, number).Scan(&ph.Id, &ph.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &ph, nil
}

func updateRecord(db *sql.DB, ph *Phone) error {
	statement := `UPDATE phone_numbers SET number =$2 WHERE ID=$1`
	_, err := db.Exec(statement, ph.Id, ph.Number)
	return err
}

func deleteRecord(db *sql.DB, ph *Phone) error {
	statement := `DELETE FROM phone_numbers WHERE ID=$1`
	_, err := db.Exec(statement, ph.Id)
	return err
}
