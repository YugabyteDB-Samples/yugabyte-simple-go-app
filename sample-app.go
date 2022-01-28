package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host        = "2e82c651-1d80-47f2-8cdf-07612b4fbda8.aws.ybdb.io"
	port        = 5433
	dbName      = "yugabyte"
	dbUser      = "admin"
	dbPassword  = "DDkxo5ZhxzPgGTdxYrGfq7Ib-FRI0o"
	sslMode     = "verify-full"
	sslRootCert = "/Users/dmagda/Downloads/yb_cloud/root.crt"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s",
		host, port, dbUser, dbPassword, dbName, sslMode, sslRootCert)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Read from the table.
	var return_val int
	rows, err := db.Query(`SELECT 1`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Printf("Query returned: ")
	for rows.Next() {
		err := rows.Scan(&return_val)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Row[%d]\n", return_val)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}
