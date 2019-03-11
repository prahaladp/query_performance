package main
import (
    "database/sql"
    "time"
    "fmt"
)

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

var queryHostMap = map[string]QueryHostParams{}

func makeNewQueryHostMap() {
    queryHostMap = map[string]QueryHostParams{}
}

func newqueryHostParams(hName string) QueryHostParams {
    var nqp = QueryHostParams{hostName:hName}
    logger.Println(" adding ", nqp)
    return nqp
}

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

func getqueryHostParams(hostName string) (QueryHostParams, bool) {
    qp, exists := queryHostMap[hostName]
    logger.Println(" query for ", hostName, " returned ", qp, exists)
    return qp, exists
}

func addOneMinute(t time.Time) time.Time {
    return t.Add(time.Minute)
}

func generateMinMaxQueryString(table string, host string,
    startTs time.Time, endTs time.Time) string {
    queryStr := fmt.Sprintf("SELECT time_bucket('1 minute', ts) AS one_min," +
        "max(usage), min(usage) from %s where host='%s' " +
        "and ts >='%s' and ts <= '%s' GROUP BY one_min",
        table, host, convertTimeToString(startTs),
        convertTimeToString(endTs))

    return queryStr
}

func extractTimeAndStore(r *sql.Rows) []TimeStampResults {
    var minU    float32
    var maxU    float32
    var t       time.Time

    usageD := []TimeStampResults{}

    if r == nil {
        return usageD
    }

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

func (qhr *QueryHostParams) findMinMaxForInterval(db Database, tr TimeRange,
        qr *QueryHostResults) {
    queryStr := generateMinMaxQueryString(db.getTableName(),
        qhr.hostName, tr.st, tr.et)

    usageData, d, e := db.queryWithTime(queryStr)
    qr.timeTaken = append(qr.timeTaken, int64(d))
    if e != nil {
        logger.Printf("query %s failed with %s", queryStr, e)
        return
    }
    logger.Printf("number of entries from database query = %d\n", len(usageData))
    logger.Printf("existing entries in usage = %d\n", len(qr.usageResults))
    qr.usageResults = append(qr.usageResults, usageData...)
    logger.Printf("new number of entries = %d\n", len(qr.usageResults))
}

func (qhr *QueryHostParams) findMinMaxForIntervals(db Database) *QueryHostResults {
    qResults := QueryHostResults{hostName: qhr.hostName}

    for _, trange := range qhr.trList {
        qhr.findMinMaxForInterval(db, trange, &qResults)
    }
    logger.Println("queryHostParams :")
    logger.Println(qhr)
    logger.Println("queryHostParams ----")
    return &qResults
}

func printHostResults(results []*QueryHostResults) {
    for _, q := range results {
        fmt.Printf("Hostname : %s\n", q.hostName)  
        for _, ur := range q.usageResults {
            fmt.Printf("\t%s, %2.2f, %2.2f\n", convertTimeToString(ur.ts),
                ur.mn, ur.mx)
        }
    }
}
