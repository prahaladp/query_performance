package main

import (
    "database/sql"
    "fmt"
    "time"

  _ "github.com/lib/pq"
)

type dbInst struct {
    host        string
    port        int
    user        string
    password    string
    tablename   string
    dbConn      *sql.DB
}

func createNewDbInst(host string, port int, user string,
    password string, table string) (dbInst, error) { 
    var db = dbInst{host: host, port: port, user: user,
        password: password, tablename: table}
    err := db.connect() 
    if err != nil {
        fmt.Println("db connect failed ", err)
        logger.Println("db connect failed ", err)
    }
    return db, err
}

func (db *dbInst) connect() error {
    connStr := "user=" + db.user + " password=" + db.password +
        " sslmode=disable"
    var err error
    db.dbConn, err = sql.Open("postgres", connStr)
    return err
}

func (db *dbInst) query(qStr string) (*sql.Rows, error) {
    rows, err := db.dbConn.Query(qStr)
    if err != nil {
        fmt.Printf("db query failed %s\n", err)
    }
    return rows, err
}

func (db dbInst) queryWithTime(qStr string) ([]TimeStampResults, time.Duration, error) {
    logger.Println("executing query : ", qStr)
    beginT := time.Now()
    rows, err := db.query(qStr)
    durT := time.Since(beginT)
    return extractTimeAndStore(rows), durT, err
}

func (db *dbInst) exec(execStr string) error {
    _, err := db.dbConn.Exec(execStr)
    if err != nil {
        fmt.Println("db exec for " + execStr)
        logger.Println("db exec for " + execStr)
        logger.Println(err)
    }
    return err
}

func (db *dbInst) clearTable() error {
    err := db.exec("DELETE from " + db.tablename) 
    return err
}

func (db dbInst) getTableName() string {
    return db.tablename
}

func (db *dbInst) close() {
    db.dbConn.Close()
}
