package main

import (
    "testing"
    "os"
)

func TestMain(m *testing.M) {
    initLog()
    os.Exit(m.Run())
}

func TestBasicFile(t *testing.T) {
    err := parseCsv("./testdata/qp.csv")
    if err != nil {
        t.Errorf("Unable to open file")
    }
}
