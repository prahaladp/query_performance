package main

import (
   "log"
   "os"
   "flag"
    "sync"
    "fmt"

  _ "github.com/lib/pq"
)

// set of variables from the user for the command line tool
// NOTE : 
// (1) see usage
// (2) can be potentially improved to load the configuration from a file
var (
  host      *string
  port      *int
  user      *string 
  password  *string
  tablename *string 
  filename  *string
  workers   *int
)

// logger
var logger  *log.Logger
var f       *os.File
var err     error

const (
    LOGFILE = "run.log"
)

// Initializes the log file
func initLog() {
    f, err = os.OpenFile(LOGFILE,
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
     }

     logger = log.New(f, "", log.LstdFlags)
     logger.Println("starting logging")
}

// main process function to synchronize across worker threads
func process() {
    var wg sync.WaitGroup

    wg.Add(*workers)

    // create channels for send and receive
    // send : the query host parameters which have the time
    //          slices to be queried for
    // recv : the set of results from the database for each query
    qParams := make(chan QueryHostParams,  100)
    results := make(chan *QueryHostResults, 100)
    dbConn := []dbInst{}

    // create the database instance and feed it to the worker
    // threads
    for i := 0; i < *workers; i++ {
        var db  dbInst

        logger.Printf("%d -> creating db\n", i)
        db, err = createNewDbInst(*host, *port, *user, *password,
            *tablename)
        dbConn = append(dbConn, db)

        if err != nil {
            fmt.Println("creating a new db connection failed ")
            logger.Println("creating a new db connection failed")
            logger.Println(err)
            os.Exit(1)
        }
        defer dbConn[i].close()
    }

    // start the workers
    for i := 0; i < *workers; i++ {
        logger.Printf("setting up worker : %d\n", i)
        go queryWorker(i, dbConn[i], qParams, results, &wg)
    }

    // pass in the query parameters which have the time slices
    for _, v := range queryHostMap {
        qParams <- v
    }

    logger.Printf("Closing query hosts params channel\n")
    close(qParams)

    // wait till all the threads have completed
    wg.Wait()

    // close the result channel
    close(results)

    logger.Printf("results")


    // gather the results
    allTimes := []int64{}
    qhr := []*QueryHostResults{}
    for q := range results {
        qhr = append(qhr, q)
        allTimes = append(allTimes, q.timeTaken...)
    }

    // compute the median/mean query performance and per-host
    // statistics
    computeAndPrint(allTimes)
    printHostResults(qhr)
}

func printArgs() {
    logger.Printf("dbport = %d\n", *port)
    logger.Printf("host = %s\n", *host)
    logger.Printf("user = %s\n", *user);
    logger.Printf("password = %s\n", *password)
    logger.Printf("table = %s\n", *tablename)
    logger.Printf("filename = %s\n", *filename)
    logger.Printf("workerd = %d\n", workers)
}

func main() {
    initLog()
    defer f.Close()
    defer logger.Println("end query_performance")

    // set up flags for input parameters
    port = flag.Int("dbport", 5432, "db port to connect to")
    user = flag.String("username", "postgres", "db user name")
    password = flag.String("password", "password", "db user password")
    tablename = flag.String("tablename", "cpu_usage", "table to query")
    host = flag.String("host", "localhost", "destination hostname for db")
    workers = flag.Int("workers", 1, "number of worker threads")
    filename = flag.String("file", "", "filename containing the query parameters")

    flag.Parse()
    printArgs()

    if *filename == "" {
        flag.Usage()
        os.Exit(1)
    }

    // parse the input CSV file - parse would populate the
    // global QueryHostParams which contain the time slices
    err = parseCsv(*filename)
    if err != nil {
        os.Exit(1)
    }

    process()
}
