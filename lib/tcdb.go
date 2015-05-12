package lib

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Exec(query string, args ...interface{}) {
	db := getDB()
	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare(query) // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates
	//fmt.Printf("args %v\n", args)
	_, err = stmtIns.Exec(args...) // 执行插入
	if err != nil {
		panic(err.Error())
	}

}

func Login(userid, password, userType string) bool {
	db := getDB()
	defer db.Close()

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("select user_id from user where user_id =? and password=? and user_type=? and status=1")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	var uid string // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(userid, password, userType).Scan(&uid) // WHERE number = 13
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

func GetAnswer(fileMd5 string) string {
	//return ""
	db, err := sql.Open("mysql", TCConfig.DBConnectStr)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT answer from exam where  file_hash =? and answer_result =1 and answer is not null ")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	var answer string // we "scan" the result in here

	// Query the square-number of 13
	rows, err := stmtOut.Query(fileMd5)

	defer rows.Close()
	if err != nil {
		return ""
	}
	for rows.Next() {
		err = rows.Scan(&answer) // WHERE number = 13
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return answer
	}
	return answer
}

///获取系统最后开启时间
func GetSysLastStartTime() string {
	db := getDB()
	defer db.Close()

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("select active_time from user_activities where user_type='system' and active_type ='start' order by active_time desc limit 1")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	var activeTime string // we "scan" the result in here

	// Query the square-number of 13
	rows, err := stmtOut.Query()

	defer rows.Close()
	if err != nil {
		return ""
	}
	for rows.Next() {
		err = rows.Scan(&activeTime) // WHERE number = 13
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return activeTime
	}
	return activeTime
}

///根据管理员id 获取用户列表
func GetUsersByManagerId(mId string) []string {
	db := getDB()
	defer db.Close()

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("select u_id from user_group where m_id=? ")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	// Query the square-number of 13
	rows, err := stmtOut.Query(mId)

	defer rows.Close()
	if err != nil {
		return nil
	}
	ret := make([]string, 0)
	for rows.Next() {
		var uId string        // we "scan" the result in here
		err = rows.Scan(&uId) // WHERE number = 13
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		ret = append(ret, uId)
	}
	return ret
}

func QueryInt(query string, args ...interface{}) int {
	db := getDB()
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

func getDB() *sql.DB {
	db, err := sql.Open("mysql", TCConfig.DBConnectStr)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	return db
}
