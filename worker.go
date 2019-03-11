package main

import (
    "sync"
)

func queryWorker(id int, db dbInst,
    qpChan <- chan QueryHostParams, results chan<- *QueryHostResults,
    wg *sync.WaitGroup) {
    for {
        qp, more := <-qpChan
        if more {
            logger.Printf("worker %d : hostname %s", id, qp.hostName)
            qhr := qp.findMinMaxForIntervals(db)
            logger.Printf("worker %d: queryresults \n", id)
            logger.Println(qhr)
            results <- qhr
        } else {
            wg.Done()
            logger.Println("query host parameters channel closed")
            return
        }
    }
}

