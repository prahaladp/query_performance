package main

import (
    "fmt"
    "sync"
)

func queryWorker(db dbInst,
    qpChan <- chan QueryHostParams, results chan<- *QueryHostResults,
    wg *sync.WaitGroup) {
    for {
        qp, more := <-qpChan
        if more {
            fmt.Println("hostname : " + qp.hostName)
            qhr := qp.findMinMaxForIntervals(db)
            fmt.Printf("worker : queryresults \n")
            fmt.Println(qhr)
            results <- qhr
        } else {
            wg.Done()
            close(results)
            fmt.Println("query host parameters channel closed")
            return
        }
    }
}

