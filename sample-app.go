package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s"+
		" dbname=%s sslmode=%s sslrootcert=%s",
		host, port, dbUser, dbPassword, dbName, sslMode, sslRootCert)

	db, err := sql.Open("postgres", psqlInfo)
	checkIfError(err)

	defer db.Close()

	fmt.Println(">>>> Successfully connected to YugabyteDB!")

	createDatabase(db)
	selectAccounts(db)
	transferMoneyBetweenAccount(db, 800)
	selectAccounts(db)
}

func createDatabase(db *sql.DB) {
	stmt := `DROP TABLE IF EXISTS DemoAccount`
	_, err := db.Exec(stmt)
	checkIfError(err)

	stmt = `CREATE TABLE IF NOT EXISTS DemoAccount (
						id int PRIMARY KEY,
						name varchar,
						age int,
						country varchar,
						balance int)`

	_, err = db.Exec(stmt)
	checkIfError(err)

	stmt = `INSERT INTO DemoAccount VALUES 
				(1, 'Jessica', 28, 'USA', 10000),
				(2, 'John', 28, 'Canada', 9000)`

	_, err = db.Exec(stmt)
	checkIfError(err)

	fmt.Println(">>>> Successfully created table DemoAccount.")
}

func selectAccounts(db *sql.DB) {
	fmt.Println(">>>> Selecting accounts:")

	rows, err := db.Query("SELECT name, age, country, balance FROM DemoAccount")
	checkIfError(err)

	defer rows.Close()

	var name, country string
	var age, balance int

	for rows.Next() {
		err = rows.Scan(&name, &age, &country, &balance)
		checkIfError(err)

		fmt.Printf("name = %s, age = %v, country = %s, balance = %v\n",
			name, age, country, balance)
	}
}

func transferMoneyBetweenAccount(db *sql.DB, amount int) {
	tx, err := db.Begin()
	checkIfError(err)

	_, err = tx.Exec(`UPDATE DemoAccount SET balance = balance - $1 WHERE name = 'Jessica'`, amount)
	if checkIfTxAborted(err) {
		return
	}
	_, err = tx.Exec(`UPDATE DemoAccount SET balance = balance + $1 WHERE name = 'John'`, amount)
	if checkIfTxAborted(err) {
		return
	}

	err = tx.Commit()
	if checkIfTxAborted(err) {
		return
	}

	fmt.Printf(">>>> Transferred %v between accounts.\n", amount)
}

func checkIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkIfTxAborted(err error) bool {
	if err != nil {
		pqErr := err.(*pq.Error)

		if pqErr.Code == `40001` {
			fmt.Println(
				`The operation is aborted due to a concurrent transaction that is modifying the same set of rows.
				 Consider adding retry logic for production-grade applications.`)
			return true

		}

		log.Fatal(err)
	}

	return false
}
