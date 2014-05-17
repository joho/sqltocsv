// sqltocsv is a package to make it dead easy to turn arbitrary database query
// results (in the form of database/sql Rows) into CSV output.
package sqltocsv

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// WriteCsvToFile will write a CSV file to the file name specified (with headers)
// based on whatever is in the sql.Rows you pass in. It calls WriteCsvToWriter under
// the hood.
func WriteCsvToFile(csvFileName string, rows *sql.Rows) error {
	f, err := os.Create(csvFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return WriteCsvToWriter(f, rows)
}

// WriteCsvToFile will write a CSV file to the writer passed in (with headers)
// based on whatever is in the sql.Rows you pass in.
func WriteCsvToWriter(writer io.Writer, rows *sql.Rows) error {
	csvWriter := csv.NewWriter(writer)

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	csvWriter.Write(columns)

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	row := make([]string, count)

	for rows.Next() {

		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}

		if err = rows.Scan(valuePtrs...); err != nil {
			return err
		}

		for i, _ := range columns {
			var value interface{}
			rawValue := values[i]

			byteArray, ok := rawValue.([]byte)

			if ok {
				value = string(byteArray)
			} else {
				value = rawValue
			}

			row[i] = fmt.Sprintf("%v", value)
		}
		csvWriter.Write(row)
	}
	err = rows.Err()

	csvWriter.Flush()

	return err
}
