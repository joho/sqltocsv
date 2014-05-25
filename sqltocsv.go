// sqltocsv is a package to make it dead easy to turn arbitrary database query
// results (in the form of database/sql Rows) into CSV output.
package sqltocsv

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// WriteFile will write a CSV file to the file name specified (with headers)
// based on whatever is in the sql.Rows you pass in. It calls WriteCsvToWriter under
// the hood.
func WriteFile(csvFileName string, rows *sql.Rows) error {
	f, err := os.Create(csvFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return Write(f, rows)
}

// WriteString will return a string of the CSV. Don't use this unless you've
// got a small data set or a lot of memory
func WriteString(rows *sql.Rows) (string, error) {
	buffer := &bytes.Buffer{}

	err := Write(buffer, rows)

	if err != nil {
		return "", err
	} else {
		return buffer.String(), nil
	}
}

// Write will write a CSV file to the writer passed in (with headers)
// based on whatever is in the sql.Rows you pass in.
func Write(writer io.Writer, rows *sql.Rows) error {
	return New(rows).Write(writer)
}

// Converter does the actual work of converting the rows to CSV.
// There are a few settings you can override if you want to do
// some fancy stuff to your CSV.
type Converter struct {
	rows    *sql.Rows
	Headers []string
}

// String returns the CSV as a string in an fmt package friendly way
func (c Converter) String() string {
	csv, err := WriteString(c.rows)
	if err != nil {
		return ""
	} else {
		return csv
	}
}

// Write writes the CSV to the Writer provided
func (c Converter) Write(writer io.Writer) error {
	rows := c.rows
	csvWriter := csv.NewWriter(writer)

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// use Headers if set, otherwise default to
	// query Columns
	var headers []string
	// TODO remove when I've figured out what's going on
	fmt.Printf("%v", c.Headers)
	if len(c.Headers) > 0 {
		headers = c.Headers
	} else {
		headers = columns
	}
	csvWriter.Write(headers)

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

// New will return a Converter which will write your CSV however you like
// but will allow you to set a bunch of non-default behaivour like overriding
// headers or injecting a pre-processing step into your conversion
func New(rows *sql.Rows) *Converter {
	return &Converter{
		rows: rows,
	}
}
