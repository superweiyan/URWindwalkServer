package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlDB struct {
	db *sql.DB
}

// func main() {
// 	db, err := sql.Open("mysql", "weiyan:Lin123456789@tcp(localhost:3306)/urwindwalk?charset=utf8")

// 	if err != nil {
// 		fmt.Println("database initialize error")
// 		panic(err.Error())
// 	}

// 	defer db.Close()

// 	mysqlDBObj := new(mysqlDB)
// 	mysqlDBObj.db = db

// 	mysqlDBObj.insertData("INSERT INTO registertable( passport, password, registerdate ) VALUES( ?, ?, ? )", "weiyan", "12345", getTimeStamp())

// 	rows, err := db.Query("select * from registertable")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer rows.Close()

// 	json := mysqlDBObj.selectData("select * from registertable")

// 	fmt.Printf(json)
// }

func (mDB *mysqlDB) updateDate(sql string) {
	if mDB.db == nil {
		return
	}

	stmt, err := mDB.db.Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()
}

func (mDB *mysqlDB) insertData(sql string, args ...interface{}) bool {
	if mDB.db == nil {
		fmt.Printf("db is nil")
		return false
	}

	stmtIns, err := mDB.db.Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
		// panic(err.Error())
		return false
	}

	defer stmtIns.Close()

	result, err := stmtIns.Exec(args...)

	if err != nil {
		// panic(err.Error())
		fmt.Println(err.Error())
		return false
	}

	if id, err := result.LastInsertId(); err == nil {
		fmt.Println("insert id : ", id)
	}

	return true
}

func (mDB *mysqlDB) selectData(sqlString string) string {
	if mDB.db == nil {
		return ""
	}

	stmt, err := mDB.db.Prepare(sqlString)
	if err != nil {
		fmt.Println("Query Error", err)
		return "Error"
	}

	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("Query Error", err)
		return "Error"
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for index := 0; index < len(columns); index++ {
		fmt.Printf("columns[%d]=%s \n", index, columns[index])
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	var jsonstring string
	jsonstring = "{\"data\":["
	allcount := 0

	for rows.Next() {
		jsonstring += "{"
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		allcount := 0
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//          fmt.Println(columns[i], ": ", value)
			if i == len(values)-1 {
				jsonstring += "\"" + columns[i] + "\":\"" + value + "\""
			} else {
				jsonstring += "\"" + columns[i] + "\":\"" + value + "\","
			}
			//          fmt.Println(" :", i, ": ", col, len(values))
		}

		jsonstring += "},"
		allcount++
	}
	if allcount > 0 {
		fmt.Println(",,,," + jsonstring)
		jsonstring = substring(jsonstring, 0, len(jsonstring)-2)
	}
	jsonstring += "]}"

	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return jsonstring
}

func getTimeStamp() int64 {
	now := time.Now().Unix()
	return now
}

func substring(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
