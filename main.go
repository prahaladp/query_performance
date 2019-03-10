package main

import (
  "database/sql"
  "fmt"
   "log"

  _ "github.com/lib/pq"
)

const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = "password"
  dbname   = "cpu_usage"
)

func main() {
  connStr := "user=postgres password=password sslmode=disable"
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    panic(err)
  }
  defer db.Close()

    fmt.Println("enter..")
  host := " host1"
  rows, err := db.Query(`SELECT min(usage) from cpu_usage where host=$1`, host)
var (
    usage float32
)
  defer rows.Close()
  for rows.Next() {
    err := rows.Scan(&usage)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(usage)
  }
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  fmt.Println("Successfully connected!")
}
