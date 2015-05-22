package mysql

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func (m *MySQL) collect() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic: MySQL: %v\n", r)
			m.clear()
		}
	}()

	db, err := sql.Open("mysql", m.cs)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("SHOW STATUS")
	if err != nil {
		panic(err.Error())
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		value, _ := strconv.Atoi(string(values[1]))

		switch string(values[0]) {
		case "Queries":
			m.values[queriesGauge].set(value)
		case "Slow_queries":
			m.values[slowGauge].set(value)
		}
	}

}
