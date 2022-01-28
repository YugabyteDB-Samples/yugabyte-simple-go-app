package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host        = ""
	port        = 5433
	dbName      = "yugabyte"
	dbUser      = ""
	dbPassword  = ""
	sslMode     = ""
	sslRootCert = ""
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s",
		host, port, dbUser, dbPassword, dbName, sslMode, sslRootCert)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(">>>> Successfully connected to YugabyteDB!")

	createDatabase(db)
	selectAccounts(db)

	defer db.Close()
}

func createDatabase(db *sql.DB) {
	stmt := `CREATE TABLE IF NOT EXISTS DemoAccount (
						id int PRIMARY KEY,
						name varchar,
						age int,
						country varchar,
						balance int)`

	_, err := db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	stmt = `INSERT INTO DemoAccount VALUES 
				(1, 'Jessica', 28, 'USA', 10000),
				(2, 'John', 28, 'Canada', 9000)`

	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(">>>> Successfully created DemoAccount table.")
}

func selectAccounts(db *sql.DB) {
	fmt.Println(">>>> Selecting accounts:")

	rows, err := db.Query("SELECT name, age, country, balance FROM DemoAccount")
	if err != nil {
		log.Fatal(err)
	}

	var name, country string
	var age, balance int

	for rows.Next() {
		err = rows.Scan(&name, &age, &country, &balance)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("name = %s, age = %v, country = %s, balance = %v\n", name, age, country, balance)
	}

	defer rows.Close()
}
