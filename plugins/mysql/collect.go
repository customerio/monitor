package mysql

import (
	"database/sql"
	"strconv"

	"github.com/customerio/monitor/metrics"
	"github.com/customerio/monitor/plugins"
	_ "github.com/go-sql-driver/mysql"
)

var updaterMap = map[string]int{
	"Queries":      queriesGauge,
	"Slow_queries": slowGauge,
}

func (f *MySQL) Collect(b *metrics.Batch) {
	f.collect()

	for _, u := range f.updaters {
		u.Fill(b)
	}
}

func (m *MySQL) collect() {
	defer func() {
		if r := recover(); r != nil {
			plugins.Logger.Printf("panic: MySQL: %v\n", r)
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

		if key, ok := updaterMap[string(values[0])]; ok {
			m.updaters[key].Update(float64(value))
		}
	}
}
