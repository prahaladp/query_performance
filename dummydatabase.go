package main

import (
    "database/sql"
    "time"
)

type DummyDatabase struct {
}

func (db DummyDatabase) queryWithTime(qStr string) (*sql.Rows, time.Duration, error) {
    beginT := time.Now()
    durT := time.Since(beginT)
    return nil, durT, nil
}

func (db DummyDatabase) getTableName() string {
    return ""
}
