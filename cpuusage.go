package main

import (
    "time"
)

// stores the cpu usage from the CSV file - this primarily
// used as a test utility function

type CpuUsageTable struct {
    ts      time.Time
    host    string
    usage   float32
}

