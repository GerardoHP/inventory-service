package data

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// )

// var DbConn *sql.DB

// func SetupDatabase() {
// 	var err error
// 	config := NewConfigFromFile()

// 	connString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", config.User, config.Pass, config.Host, config.Port, config.DBName)
// 	fmt.Printf(" conn string %v \n", connString)
// 	DbConn, err = sql.Open("mysql", connString)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
