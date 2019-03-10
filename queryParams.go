package main
import (
    "fmt"
    "database/sql"
    "time"
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
    fmt.Println(" adding ", nqp)
    return nqp
}

func addTimeRange(hostName string, tr TimeRange) {
    if qp, exists := queryHostMap[hostName]; exists == false {
        fmt.Println("Adding a new hostName ", hostName)
        qp = newqueryHostParams(hostName)
        qp.trList = append(qp.trList, tr)
        fmt.Println(" new ", qp)
        queryHostMap[hostName] = qp
    } else {
        fmt.Println("Appending to existing host")
        qp.trList = append(qp.trList, tr)
        queryHostMap[hostName] = qp
        fmt.Println(" new ", qp)
    }
}

func getqueryHostParams(hostName string) (QueryHostParams, bool) {
    qp, exists := queryHostMap[hostName]
    fmt.Println(" query for ", hostName, " returned ", qp, exists)
    return qp, exists
}

func addOneMinute(t time.Time) time.Time {
    return t.Add(time.Minute)
}

func generateMinMaxQueryString(table string, host string,
    ts time.Time) (string, string) {
    minStr := fmt.Sprintf("SELECT min(usage) from %s where host =" + 
        "'%s' and ts='%s'", table, host, convertTimeToString(ts))
    maxStr := fmt.Sprintf("SELECT max(usage) from %s where host =" +
        "'%s' and ts='%s'", table, host, convertTimeToString(ts))

    return minStr, maxStr
}

func getTimeFromRows(r *sql.Rows) float32 {
    var t float32
    t = 0.0
    if r == nil {
        return t
    }

    for r.Next() {
        err := r.Scan(&t)
        fmt.Printf("scanning %f\n", t)
        if err != nil {
            fmt.Println("get min/max usage time failed ")
            fmt.Println(err)
        }
        break
    }
    return t
}

func (qhr *QueryHostParams) findMinMaxForInterval(db Database, tr TimeRange,
        qr *QueryHostResults) {
    currT := tr.st
    endT := tr.et
    for currT.Before(endT) || currT.Equal(endT) {
        var minR, maxR *sql.Rows
        var d time.Duration
        var e error
        var minTime, maxTime float32

        minStr, maxStr := generateMinMaxQueryString(
            db.getTableName(), qhr.hostName,
            currT)
        currT = addOneMinute(currT)

        // NOTE : we count query times even if the query failed or 
        // if the query returned no results

        minR, d, e = db.queryWithTime(minStr)

        qr.timeTaken = append(qr.timeTaken, int64(d))
        if minR != nil { 
            minTime = getTimeFromRows(minR)
            defer minR.Close()
        }        

        if e != nil {
            fmt.Printf("query %s failed with %s", minStr, e)
            continue
        }

        maxR, d, e = db.queryWithTime(maxStr)

        qr.timeTaken = append(qr.timeTaken, int64(d))

        if maxR != nil {
            maxTime = getTimeFromRows(maxR)
            defer maxR.Close()
        }

        if e != nil {
            fmt.Printf("query %s failed with %s", maxStr, e)
            continue
        }

        usageData := TimeStampResults{ts: currT, mn: minTime, mx: maxTime}
        qr.usageResults = append(qr.usageResults, usageData)
    }
}

func (qhr *QueryHostParams) findMinMaxForIntervals(db Database) *QueryHostResults {
    qResults := QueryHostResults{hostName: qhr.hostName}

    for _, trange := range qhr.trList {
        qhr.findMinMaxForInterval(db, trange, &qResults)
    }
    fmt.Println("queryHostParams :")
    fmt.Println(qhr)
    fmt.Println("queryHostParams ----")
    return &qResults
}
