package main

import (
    "sync"
)

// Main worker thread.
// The workers are identified by an id and the caller maintains
// the number of workers.
// There are two channel - each host's QueryHostParams is read from
// the qpChan (channel)
// The output is sent back on the results channel (type QueryHostResults)
// The worker keeps processing till the qpChan is closed by the
// the caller
// At the end the WaitGroup is signalled indicating that the worker
// is done processing
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

