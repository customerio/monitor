package mysql

import (
	"database/sql"
	"strconv"
)

func (m *MySQL) collect() {

	db, err := sql.Open("mysql", m.cs)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
	    panic(err.Error()) // proper error handling instead of panic in your app
	}

	rows, err := db.Query("SHOW STATUS")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }


    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
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
            panic(err.Error()) // proper error handling instead of panic in your app
        }	

        value, _ := strconv.Atoi(string(values[1]))

        switch string(values[0]){
        	case "Queries": m.queries.Set(value)
        	case "Slow_queries": m.slow.Set(value)
        }
    }

}
