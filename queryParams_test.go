package main

import (
    "testing"
    "time"
)

func TestBasic(t *testing.T) {
    makeNewQueryHostMap()    
    addTimeRange("h1", TimeRange{time.Now(), time.Now()})        
    addTimeRange("h2", TimeRange{time.Now(), time.Now()})
    if len(queryHostMap) != 2 {
        t.Errorf("Number of map elements is not 2")
    }
    qp, exists := getqueryHostParams("h1")
    if !exists {
        t.Errorf("Didn't find host h1");
    }
    if len(qp.trList) != 1 {
        t.Errorf("Number of elements in time range is incorrect")
    }

    qp, exists = getqueryHostParams("bbbfds")
    if exists  {
        t.Errorf("Expecting to not find host bbbfds")
    }
}

func TestQueryResults1(t *testing.T) {
    makeNewQueryHostMap()
    
    st := time.Now()
    et := st.Add(time.Minute * 10)
    addTimeRange("h1", TimeRange{st, et})

    st = st.Add(time.Hour * 1)
    et = st.Add(time.Hour * 1)
    addTimeRange("h1", TimeRange{st, et})

    qp, exists := getqueryHostParams("h1")
    if !exists {
        t.Errorf("Didn't find host h1");
    }

    logger.Printf("Number of entries = %d\n", len(qp.trList))
    if len(qp.trList) != 2 {
        t.Errorf("Number of elements in time range is incorrect")
    }

    db := DummyDatabase{}
    usageResults := []TimeStampResults{}
    usageResults = append(usageResults, TimeStampResults{ts:st.Add(time.Minute),
        mn: 50.0, mx: 67.1})
    usageResults = append(usageResults, TimeStampResults{ts:st.Add(time.Minute * 3),
            mn: 50.0, mx: 67.1})
    db.setTimeStampResults(usageResults)
    qhr := qp.findMinMaxForIntervals(db)

    logger.Printf("Usage results = %d\n", len(qhr.usageResults))
    logger.Printf("Time samples = %d\n", len(qhr.timeTaken))

    if len(qhr.usageResults) != 4 {
        t.Errorf("Incorrect number of usage results")
    }
}
