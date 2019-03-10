package main

import (
    "time"
    "database/sql"
)

type Database interface {  
    queryWithTime(string) (*sql.Rows, time.Duration, error)
    getTableName() string
}

