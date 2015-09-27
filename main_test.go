package main_test

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func Benchmark_ids(b *testing.B) {
	db, err := sql.Open("mysql", "root:secret1234@tcp(192.168.33.10:3306)/")
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec(`CREATE DATABASE IF NOT EXISTS ids`)
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec(`USE ids`)
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec(`DROP TABLE id`)
	if err != nil {
		b.Fatal(err)
	}

	const createID = `CREATE TABLE IF NOT EXISTS id (
  id INT(32),
  PRIMARY KEY (id)
) ENGINE=innodb`

	_, err = db.Exec(createID)
	if err != nil {
		b.Fatal(err)
	}

	statement, err := db.Prepare("INSERT INTO id VALUES (?)")
	if err != nil {
		b.Fatal("Boom prepared statement creation failed.")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := statement.Exec(i)
		if err != nil {
			b.Fatal(err)
		}
	}

}

func Benchmark_NewUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUUID()
	}
}

func Benchmark_uuids(b *testing.B) {
	db, err := sql.Open("mysql", "root:secret1234@tcp(192.168.33.10:3306)/")
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec(`CREATE DATABASE IF NOT EXISTS uids`)
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec(`USE uids`)
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.Exec(`DROP TABLE uuid`)
	if err != nil {
		b.Fatal(err)
	}

	const createUUID = `CREATE TABLE IF NOT EXISTS uuid (
  id BINARY(16),
  PRIMARY KEY (id)
) ENGINE=innodb`

	_, err = db.Exec(createUUID)
	if err != nil {
		b.Fatal(err)
	}

	uuids := make([][]byte, b.N)

	for uindex := 0; uindex < b.N; uindex++ {
		uuids[uindex] = NewUUID()
	}
	// Reset the time as we don't want to time the UUID creation process

	statement, err := db.Prepare("INSERT INTO uuid VALUES (?)")
	if err != nil {
		b.Fatal("Boom prepared statement creation failed.")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := statement.Exec(uuids[i])
		if err != nil {
			b.Fatal(err)
		}
	}
}
