package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// testing procedure
// - insert records serially tracking response time for each.
// - draw a distribution graph of the throughput.
//

func NewUUID() []byte {
	b := make([]byte, 16)

	rand.Read(b)

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return b
}

type LoopPtr func(*sql.DB, int)

func main() {
	var insertCount int
	var idType string

	flag.IntVar(&insertCount, "insert", 10000, "Number of records to insert.")
	flag.StringVar(&idType, "type", "id", "Type of id to use (id, uuid).")
	flag.Parse()

	db, err := sql.Open("mysql", "root:secret1234@tcp(192.168.33.10:3306)/")
	if err != nil {
		log.Fatal(err)
	}

	var fn LoopPtr
	var tableStatement string
	switch idType {
	case "id":
		fn = IdInsert
		tableStatement = `CREATE TEMPORARY TABLE id (
  id INT(32),
  PRIMARY KEY (id)
) ENGINE=innodb`

	case "uuid":
		fn = UuidInsert
		tableStatement = `CREATE TEMPORARY TABLE uuid (
  id BINARY(16),
  PRIMARY KEY (id)
) ENGINE=innodb`
	}

	PrepareDatabase(db, tableStatement)
	fn(db, insertCount)
}

type TrackerEntry struct {
	Time     time.Time
	Duration time.Duration
}

func (te *TrackerEntry) StrTime() string {
	return te.Time.Format(time.RFC3339)
}

func (te *TrackerEntry) StrDuration() string {
	return strconv.FormatInt(int64(te.Duration/time.Millisecond), 10)
}

type Tracker struct {
	Entries []TrackerEntry
}

func NewTracker(insertCount int) *Tracker {
	return &Tracker{
		Entries: make([]TrackerEntry, insertCount, insertCount),
	}
}

func (tracker *Tracker) Record(start time.Time, i int) {
	duration := time.Now().Sub(start)
	tracker.Entries[i] = TrackerEntry{
		Time:     start,
		Duration: duration,
	}
}

func (tracker *Tracker) Save(w io.Writer) error {
	cw := csv.NewWriter(w)
	cw.Write([]string{"time", "duration"})

	for _, entry := range tracker.Entries {
		cw.Write([]string{entry.StrTime(), entry.StrDuration()})
	}

	cw.Flush()
	return cw.Error()
}

func IdInsert(db *sql.DB, insertCount int) {
	statement, err := db.Prepare("INSERT INTO id VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}

	tracker := NewTracker(insertCount)

	start := time.Now()

	for i := 0; i < insertCount; i++ {
		t := time.Now()
		_, err := statement.Exec(i)
		tracker.Record(t, i)

		if err != nil {
			log.Fatal(err)
		}
	}

	duration := time.Now().Sub(start)
	log.Println(duration)

	file, err := os.Create("id.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = tracker.Save(file)
	if err != nil {
		log.Fatal(err)
	}
}

func UuidInsert(db *sql.DB, insertCount int) {
	statement, err := db.Prepare("INSERT INTO uuid VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}

	uuids := make([][]byte, insertCount)

	for uindex := 0; uindex < insertCount; uindex++ {
		uuids[uindex] = NewUUID()
	}

	tracker := NewTracker(insertCount)

	start := time.Now()

	for i := 0; i < insertCount; i++ {
		t := time.Now()
		_, err := statement.Exec(uuids[i])
		tracker.Record(t, i)

		if err != nil {
			log.Fatal(err)
		}
	}

	duration := time.Now().Sub(start)
	log.Println(duration)

	file, err := os.Create("uuid.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = tracker.Save(file)
	if err != nil {
		log.Fatal(err)
	}
}

func PrepareDatabase(db *sql.DB, tableStatement string) {
	_, err := db.Exec(`CREATE DATABASE IF NOT EXISTS ids`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`USE ids`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(tableStatement)
	if err != nil {
		log.Fatal(err)
	}
}
