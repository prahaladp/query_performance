package main

import (
    "testing"
    "fmt"
)

func TestMedian(t *testing.T) {
    var val []int64

    m := computeMean(val)
    if m != 0 {
        t.Errorf("Incorrect mean for 0 length array");
    }

    val = []int64{1, 2, 3}
    m = computeMedian(val)
    if m != 2 {
        t.Errorf("Incorrect median for odd size array")
        fmt.Println(val)
    }

    val = append(val, 4)
    m = computeMedian(val)
    if m != 2.5 {
        t.Errorf("Incorrect median for even size array")
        fmt.Println(m)
    }
}
