package main

import (
    "sort"
    "fmt"
    "time"
)

// set of compute utilities for calculating the mean and median
// of an array of int64's - these represent the time taken for
// the queries

type int64arr []int64
func (a int64arr) Len() int { return len(a) }
func (a int64arr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a int64arr) Less(i, j int) bool { return a[i] < a[j] }

func computeMedian(val int64arr) float64 {
    logger.Printf("size = %d\n", len(val))
    logger.Println(val)

    sort.Sort(val)
    if len(val) == 0 {
        return 0
    }

    mid := len(val) / 2
    if len(val) % 2 == 1 {
        return float64(val[mid])
    }
    return float64(val[mid-1] + val[mid])/float64(2.0)
}

func computeMean(val int64arr) float64 {
    logger.Printf("size = %d\n", len(val))
    logger.Println(val)

    var sum int64

    if len(val) == 0 {
        return 0.0
    }
    sum = 0
    for _, v := range val {
        sum += v
    }
    return float64(sum/(int64)(len(val)))
}

func computeAndPrint(allTimes []int64) {
    median := computeMedian(allTimes)
    mean := computeMean(allTimes)
    fmt.Printf("---------------------------------\n")
    fmt.Printf(" number of samples = %d\n", len(allTimes))
    fmt.Printf(" mean duration = %s\n", time.Duration(mean).String())
    fmt.Printf(" median duration = %s\n", time.Duration(median).String())
    fmt.Printf("---------------------------------\n")
}
