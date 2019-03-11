package main
import (
    "database/sql"
    "time"
    "fmt"
)

// Primary datastructures useful for the query_performance tool
// These are as follows :
//  (1) QueryHostParams - time slices for which information is required
//                        per host. The time slices are stored in the
//                        TimeRange struct
//  (2) TimeRange       - contains the start and end time
//  (3) QueryHostResults- this is used for storing the results after
//                        the queries has been issued. The durations
//                        for each query is stored in a int64 array.
//                        The min/max is stored in TimeStampResults
//  (4) TimeStampResults- The min/max for each time slot is stored
//                        in this struct

type TimeRange struct {
    st              time.Time
    et              time.Time
}

type QueryHostParams struct {
    hostName        string
    trList          []TimeRange
}

type TimeStampResults struct {
    ts              time.Time
    mx              float32
    mn              float32
}

type QueryHostResults struct {
    hostName        string
    usageResults    []TimeStampResults
    timeTaken       []int64
}

// global map indexed by host name
var queryHostMap = map[string]QueryHostParams{}

// clears the map (used during testing)
func makeNewQueryHostMap() {
    queryHostMap = map[string]QueryHostParams{}
}

func newqueryHostParams(hName string) QueryHostParams {
    var nqp = QueryHostParams{hostName:hName}
    logger.Println(" adding ", nqp)
    return nqp
}

// adds a new time range - it first looks up the map to
// see if the host name exists and adds this time slice to
// the end of the list
// NOTE : there is no optimization to concatenate the time slices
//        to reduce the number of queries - the user input is used
//        as is
func addTimeRange(hostName string, tr TimeRange) {
    if qp, exists := queryHostMap[hostName]; exists == false {
        logger.Println("Adding a new hostName ", hostName)
        qp = newqueryHostParams(hostName)
        qp.trList = append(qp.trList, tr)
        logger.Println(" new ", qp)
        queryHostMap[hostName] = qp
    } else {
        logger.Println("Appending to existing host")
        qp.trList = append(qp.trList, tr)
        queryHostMap[hostName] = qp
        logger.Println(" new ", qp)
    }
}

// looks up a specific host
func getqueryHostParams(hostName string) (QueryHostParams, bool) {
    qp, exists := queryHostMap[hostName]
    logger.Println(" query for ", hostName, " returned ", qp, exists)
    return qp, exists
}

// helper function
func addOneMinute(t time.Time) time.Time {
    return t.Add(time.Minute)
}

// generates a SQL query string using time buckets
func generateMinMaxQueryString(table string, host string,
    startTs time.Time, endTs time.Time) string {
    queryStr := fmt.Sprintf("SELECT time_bucket('1 minute', ts) AS one_min," +
        "max(usage), min(usage) from %s where host='%s' " +
        "and ts >='%s' and ts <= '%s' GROUP BY one_min",
        table, host, convertTimeToString(startTs),
        convertTimeToString(endTs))

    return queryStr
}

// extracts data from the SQL result rows and converts it into
// TimeStampResults
func extractTimeAndStore(r *sql.Rows) []TimeStampResults {
    var minU    float32
    var maxU    float32
    var t       time.Time

    usageD := []TimeStampResults{}

    if r == nil {
        return usageD
    }

    // iterate on all the rows for which data was found
    // in the database
    for r.Next() {
        err := r.Scan(&t, &maxU, &minU)
        logger.Printf("row scan generated %s %f %f\n",
            convertTimeToString(t), minU, maxU)
        if err != nil {
            logger.Println("get min/max usage time failed ")
            logger.Println(err)
        }
        usageData := TimeStampResults{ts: t, mn: minU, mx: maxU}
        usageD = append(usageD, usageData)
    }
    return usageD
}

// finds the min/max values for the time slot specified in the
// TimeRange.
func (qhr *QueryHostParams) findMinMaxForInterval(db Database, tr TimeRange,
        qr *QueryHostResults) {
    queryStr := generateMinMaxQueryString(db.getTableName(),
        qhr.hostName, tr.st, tr.et)

    // query the database - output is the SQL rows, duration and
    // the error
    usageData, d, e := db.queryWithTime(queryStr)
    qr.timeTaken = append(qr.timeTaken, int64(d))
    if e != nil {
        logger.Printf("query %s failed with %s", queryStr, e)
        return
    }

    // add the results back to QueryHostResults
    logger.Printf("number of entries from database query = %d\n", len(usageData))
    logger.Printf("existing entries in usage = %d\n", len(qr.usageResults))
    qr.usageResults = append(qr.usageResults, usageData...)
    logger.Printf("new number of entries = %d\n", len(qr.usageResults))
}

// iterate over the  input time slices by calling into
// findMinMaxForInterval
func (qhr *QueryHostParams) findMinMaxForIntervals(
    db Database) *QueryHostResults {
    qResults := QueryHostResults{hostName: qhr.hostName}

    for _, trange := range qhr.trList {
        qhr.findMinMaxForInterval(db, trange, &qResults)
    }
    logger.Println("queryHostParams :")
    logger.Println(qhr)
    logger.Println("queryHostParams ----")
    return &qResults
}

// utility function to print the data for the time slots for which
// data is available
func printHostResults(results []*QueryHostResults) {
    for _, q := range results {
        fmt.Printf("Hostname : %s\n", q.hostName)  
        for _, ur := range q.usageResults {
            fmt.Printf("\t%s, %2.2f, %2.2f\n", convertTimeToString(ur.ts),
                ur.mn, ur.mx)
        }
    }
}
