package main

import (
    "time"
)

type CpuUsageTable struct {
    ts      time.Time
    host    string
    usage   float32
}

