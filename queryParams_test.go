package main

import (
    "testing"
    "time"
    "fmt"
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

    fmt.Printf("Number of entries = %d\n", len(qp.trList))
    if len(qp.trList) != 2 {
        t.Errorf("Number of elements in time range is incorrect")
    }

    db := DummyDatabase{}
    qhr := qp.findMinMaxForIntervals(db)

    fmt.Printf("Usage results = %d\n", len(qhr.usageResults))
    if len(qhr.usageResults) != 72 {
        t.Errorf("Incorrect number of entries in the results")
    }
}
