package main

import (
    "time"
)

type Database interface {  
    queryWithTime(string) ([]TimeStampResults, time.Duration, error)
    getTableName() string
}

