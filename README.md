# sqltocsv

A library designed to let you easily turn any arbitrary sql.Rows result from a query into a CSV file with a minimum of fuss.

This is very much a work in progress at this stage.

## Usage

```go
rows, _ := db.Query("SELECT * FROM users WHERE something=72")

err := WriteCsvToFile("~/important_user_report.csv", rows)
if err != nil {
    panic(err)
}
```

## License

&copy; [John Barton](http://whoisjohnbarton.com/) but under MIT (see [LICENSE](LICENSE)) except for fakedb_test.go which I lifted from the Go standard library and is under BSD and I am unsure what that means legally.
