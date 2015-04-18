package main

import (
	"database/sql"
	"fmt"
	//	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Exec(query string, args ...interface{}) {
	db, err := sql.Open("mysql", "root:123456@tcp(10.1.9.102:3306)/proxy_exam")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare(query) // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates
	fmt.Printf("args %v\n", args)
	_, err = stmtIns.Exec(args...) // 执行插入
	if err != nil {
		panic(err.Error())
	}

}

func Login(userid, password string) bool {
	db, err := sql.Open("mysql", "root:123456@tcp(10.1.9.102:3306)/proxy_exam")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("select user_id from user where user_id =? and password=?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	var uid string // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(userid, password).Scan(&uid) // WHERE number = 13
	if err != nil {
		//panic(err.Error())
		return false
	}
	if uid == userid {
		return true
	} else {
		return false
	}
}

func QueryInt(query string, args ...interface{}) int {

	db, err := sql.Open("mysql", "root:123456@tcp(10.1.9.102:3306)/proxy_exam")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for reading data
	stmtOut, err := db.Prepare(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	var count int // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(args...).Scan(&count) // WHERE number = 13
	if err != nil {
		panic(err.Error())
	}
	return count
}
