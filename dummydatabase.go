package main

import (
    "time"
)

// represents a dummy database object for testing purposes.
type DummyDatabase struct {
    usageResults    []TimeStampResults
}

func (db DummyDatabase) queryWithTime(qStr string) (
    []TimeStampResults, time.Duration, error) {
    logger.Printf("dummydatabase : %s\n", qStr)
    logger.Println(db.usageResults)
    beginT := time.Now()
    durT := time.Since(beginT)
    return db.usageResults, durT, nil
}

func (db *DummyDatabase) setTimeStampResults(usageD []TimeStampResults) {
    db.usageResults = usageD
}

func (db DummyDatabase) getTableName() string {
    return ""
}
