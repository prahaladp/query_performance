package main

import (
    "testing"
    "fmt"
)

func TestBasicFile(t *testing.T) {
    err := parseCsv("./testdata/qp.csv")
    if err != nil {
        t.Errorf("Unable to open file")
    }
    fmt.Println(err)
}
