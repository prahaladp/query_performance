package main

import (
    "bufio"
    "encoding/csv"
    "fmt"
    // "io"
    "os"
    "time"
    "strings"
    "errors"
)

// This is used for CSV parsing and parsing the input file so
// the data is stored in QueryHostParams and subsequently used for
// SQL queries etc

// manipulates the input time string into the required format
func validateAndTransformTime(st string) (string, error) {
    var tok     []string
    var err     error

    tok = strings.Split(st, " ")
    if len(tok) != 2 {
        return st, errors.New("incorrect tokens in " + st)
    }
    return tok[0] + "T" + tok[1] + "Z", err
}

// Parses each entry in the SQL file - this is in the format
// <hostname:string>, <time:string>, <time:string>
// This calls into addTimeRange to store the parsed data in a
// map indexed by "hostname"
func parseParams(hostName string, startT string, endT string) error {
    var st, et time.Time
    var err error
    
    startT, err = validateAndTransformTime(startT)
    if err != nil {
        fmt.Println(err)
        return err
    }

    endT, err = validateAndTransformTime(endT)
    if err != nil {
        fmt.Println(err)
        return err
    }

    st, err = time.Parse(time.RFC3339, startT)
    if err != nil {
        fmt.Println(err)
        return err
    }
    et, err = time.Parse(time.RFC3339, endT)
    if err != nil {
        fmt.Println(err)
        return err
    }

    logger.Printf("Adding times %s %s\n", startT, endT)
    addTimeRange(hostName, TimeRange{st, et})
    return nil
}

// main entry point to open the file and use csv package
// to parse the file
func parseCsv(fileName string) error {
    csvFile, err := os.Open(fileName)
    if err != nil {
        fmt.Println(err)
        return err
    }

    reader := csv.NewReader(bufio.NewReader(csvFile))
    lines, e := reader.ReadAll()
    if e != nil {
        logger.Fatalf("error reading all lines: %v", err)
    }
        
    for i, line := range lines {
        if i == 0 {
            // skip the first line
            continue
        }
        logger.Printf("%s %s %s\n", line[0], line[1], line[2])
        parseParams(line[0], line[1], line[2])
    }
    return nil
}
