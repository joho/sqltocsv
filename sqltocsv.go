package sqltocsv

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
)

func WriteCsvToFile(csvFileName string, rows *sql.Rows) error {
	f, err := os.Create(csvFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	csvWriter.Write(columns)

	count := len(columns)

	for rows.Next() {
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}

		if err = rows.Scan(valuePtrs...); err != nil {
			return err
		}

		row := make([]string, count)
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
