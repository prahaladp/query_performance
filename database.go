package main

import (
    "time"
)

// Database Interface - used to distinguish between
//  (1) access to real database and
//  (2) dummy database which is used for tests
type Database interface {  
    queryWithTime(string) ([]TimeStampResults, time.Duration, error)
    getTableName() string
}

