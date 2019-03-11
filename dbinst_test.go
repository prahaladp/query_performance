package main

import (
    "fmt"
    "database/sql"
    "time"
    "testing"
)

func TestSimpleDbOps(t *testing.T) {
    const (
      host     = "localhost"
      port     = 5432
      user     = "postgres"
      password = "password"
      dbname   = "cpu_usage"
    )
    db, err := createNewDbInst(host, port, user, password,
        dbname)
    defer db.close()

    if err != nil {
        t.Errorf("creating a new instance failed ")
        logger.Println(err)
    }

    err = db.clearTable()
    if err != nil {
        t.Errorf("clearing the table failed")
        logger.Println(err)
    }

    var cUsage float32
    cUsage = 78.0
    cTimeStr := "2017-01-01 08:59:22"
    cHost := "host1"
    err = db.exec(fmt.Sprintf("INSERT into " + dbname +
        "(ts,host,usage) " +
        "VALUES (TIMESTAMP '%s'," +
        "'%s', %f)", cTimeStr, cHost, cUsage))
    if err != nil {
        t.Errorf("inserting a row failed")
        logger.Println(err)
    }

    var rows *sql.Rows
    rows, err = db.query("SELECT * from " + dbname)
    defer rows.Close()

    cpuTable := CpuUsageTable{}
    for rows.Next() {
        err = rows.Scan(&cpuTable.ts, &cpuTable.host, &cpuTable.usage)
        if err != nil {
            t.Errorf("row scan failed ");
            logger.Println(err)
        }
        logger.Println(cpuTable)
        if cpuTable.usage != cUsage {
            t.Errorf("retrieved incorrect cpu usage")
        }
        if cpuTable.host != cHost {
            t.Errorf("retrieved incorrect host name")
        }
    }
    db.clearTable()
}

func TestComponentTest1(t *testing.T) {
    makeNewQueryHostMap()

    const (
      host     = "localhost"
      port     = 5432
      user     = "postgres"
      password = "password"
      dbname   = "cpu_usage"
    )
    db, err := createNewDbInst(host, port, user, password,
        dbname)
    defer db.close()

    if err != nil {
        t.Errorf("creating a new instance failed ")
        fmt.Println(err)
    }

    err = db.clearTable()
    if err != nil {
        t.Errorf("clearing the table failed")
        fmt.Println(err)
    }

    st := time.Now()
    i := 0
    usage := []float32{1.1, 20.3, 45.1, 2.3, 6.9, 10, 98.2, 5.5, 6.2, 1.4}
    for i < 10 {
        cHost := "host1"
        cUsage := usage[i]
        // 2017-01-01 08:59:22
        cTimeStr := convertTimeToString(st)

        err = db.exec(fmt.Sprintf("INSERT into " + dbname +
            "(ts,host,usage) " +
            "VALUES (TIMESTAMP '%s'," +
            "'%s', %f)", cTimeStr, cHost, cUsage))
        if err != nil {
            t.Errorf("inserting a row failed")
            fmt.Println(err)
        }
        i++
    }

    addTimeRange("host1", TimeRange{st, st})
    qp, exists := getqueryHostParams("host1")
    if !exists {
        t.Errorf("Didn't find host host1")
    }
    
    qhr := qp.findMinMaxForIntervals(db)    
    if len(qhr.usageResults) != 1 {
        t.Errorf("Incorrect number of min/max results")
    }

    fmt.Printf("number of time samples = %d\n", len(qhr.timeTaken))
    if len(qhr.timeTaken) != 1 {
        t.Errorf("Incorrect number of time samples")
    } 
}
