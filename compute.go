package main

import (
    "sort"
)

type int64arr []int64
func (a int64arr) Len() int { return len(a) }
func (a int64arr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a int64arr) Less(i, j int) bool { return a[i] < a[j] }

func computeMedian(val int64arr) float64 {
    sort.Sort(val)
    if len(val) == 0 {
        return 0
    }

    mid := len(val) / 2
    if mid % 2 == 1 {
        return float64(val[mid])
    }
    return float64(val[mid-1] + val[mid])/float64(2.0)
}

func computeMean(val int64arr) float64 {
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
