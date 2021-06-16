package main

import (
	"fmt"

	"github.com/db"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	username = "test_user"
	password = 123456
	dbName   = "normalize_phone_numbers"
)

func checkError(err error) error {
	if err != nil {
		panic(err)
	}
	return nil
}

func main() {
	dbInfo := fmt.Sprintf("user=%s password=%d host=%s port=%d sslmode=disable", username, password, host, port)
	dbObj := db.OpenConnection("postgres", dbInfo)
	db.ResetDB(dbObj, dbName)                            //reset database and create table normalize_phone_numbers
	dbInfo = fmt.Sprintf("%s dbname=%s", dbInfo, dbName) // credential string
	dbObj = db.OpenConnection("postgres", dbInfo)        //connect to db
	db.CreateTable(dbObj)                                //create table
	db.SeedData(dbObj)                                   //insert static phone number data into table
	number_list, err := db.ViewAllRecord(dbObj)          //fetch all records from table
	checkError(err)                                      // check if any error
	db.NormalizeAndUpdateRecord(dbObj, number_list)      // remove extra spaces and delete duplicate phone numbers
	number_list, err = db.ViewAllRecord(dbObj)           // updated result with unique ph no and no spaces in between
	fmt.Println(number_list)
	// number, err = db.ViewRecord(dbObj, 1)   //returns number of specified record id
	// fmt.Println(number)	//print phone no of given record id
	defer dbObj.Close()
}
