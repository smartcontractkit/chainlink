package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DataPoint struct {
	Id        int
	Col1      int64
	Col2      int64
	KeeprInfo int64
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var batchSize = flag.Int("batch-size", 0, "batch size for break down select statements. Defaults to 0 (no batching)")
	var numRows = flag.Int("num-rows", 1000, "Number of rows to consider")
	flag.Parse()

	db := sqlx.MustConnect("postgres", `postgresql://test:test@localhost:5432/loadtest?sslmode=disable`)

	//Used to create DB
	/*
		fmt.Println("Creating Datasets.")
		schema := `CREATE TABLE IF NOT EXISTS random_data_set (
					id SERIAL,
					col1 BIGINT,
					col2 BIGINT,
					keeperId INT);`
		db.MustExec(schema)
		schema = `CREATE TABLE IF NOT EXISTS random_keepers_set (
			id SERIAL,
			keeperId INT,
			keeprInfo BIGINT);`
		db.MustExec(schema)
		fmt.Println("Datasets created.")
	*/

	//Used to insert rows
	/*
		insertRecord := `INSERT INTO random_keepers_set (keeperId, keeprInfo) VALUES ($1, $2)`
		for i := 0; i < 10; i++ {
			// random_keepers_set is small with only 10 keepers
			db.MustExec(insertRecord, i, rand.Int63())
		}

		insertRecord := `INSERT INTO random_data_set (col1, col2, keeperId) VALUES ($1, $2, $3)`
		for i := 0; i < 10000000; i++ {
			db.MustExec(insertRecord, rand.Int63(), rand.Int63(), int(rand.Int63()%10))
		}
	*/

	countSql := `SELECT count(*) FROM random_data_set`
	var total_count int

	err := db.Get(&total_count, countSql)
	if err != nil {
		fmt.Println(err)
		return

	}
	if total_count < *numRows {
		fmt.Println("Sorry not enough rows in DB\n")
		return
	}
	total_count = *numRows
	prime := int64(98447547)
	sum := int64(0)
	fmt.Printf("Going to load %d records in memory and print sum modulo %d\n", total_count, prime)

	if *batchSize == 0 {
		*batchSize = total_count
	}
	pageSize := *batchSize
	iterations := (total_count + pageSize - 1) / pageSize
	fmt.Printf("Using batch size: %d, Num Iterations: %d\n", pageSize, iterations)

	c := make(chan int64)
	for i := 0; i < iterations; i++ {
		go pageSum(db, i, pageSize, prime, c)
	}
	for i := 0; i < iterations; i++ {
		sum += <-c
		sum %= prime
	}

	fmt.Printf("Final Sum: %d\n", sum)

}

func pageSum(db *sqlx.DB, iterationNumber, pageSize int, prime int64, c chan int64) {
	sum := int64(0)
	selectSql := `SELECT random_data_set.id, col1, col2, keeprInfo FROM random_data_set INNER JOIN random_keepers_set ON random_keepers_set.keeperId = random_data_set.keeperId order by random_data_set.id asc offset $1 limit $2`
	data := []DataPoint{}
	err := db.Select(&data, selectSql, iterationNumber*pageSize, pageSize)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, dataPoint := range data {
		sum += (dataPoint.Col1 % prime)
		sum %= prime
		sum += (dataPoint.Col2 % prime)
		sum %= prime
		sum += (dataPoint.KeeprInfo % prime)
		sum %= prime
	}
	c <- sum
}
