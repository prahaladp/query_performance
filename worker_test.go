package main

import (
    "fmt"
    "time"
    "sync"

    "testing"
)

func TestSimpleWorker(t *testing.T) {
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

    // simple call
    st := time.Now()
    addTimeRange("host1", TimeRange{st, st})
    addTimeRange("host2", TimeRange{st, st.Add(time.Hour)})
    addTimeRange("host3", TimeRange{st, st.Add(time.Minute * 10)})

    // create channels
    qParams := make(chan QueryHostParams,  100)
    results := make(chan *QueryHostResults, 100)

    var wg sync.WaitGroup
    wg.Add(1)

    go queryWorker(1, db, qParams, results, &wg)

    for _, v := range queryHostMap {
        qParams <- v
    }

    logger.Printf("Closing query hosts params channel\n")
    close(qParams)
    wg.Wait()
    logger.Printf("results")

    allTimes := []int64{}
    for q := range results {
        allTimes = append(allTimes, q.timeTaken...)
    }

    computeAndPrint(allTimes)
}
