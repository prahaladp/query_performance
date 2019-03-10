package main

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "time"
    "strings"
    "errors"
)

func validateAndTransformTime(st string) (string, error) {
    var tok     []string
    var err     error

    tok = strings.Split(st, " ")
    if len(tok) != 2 {
        return st, errors.New("incorrect tokens in " + st)
    }
    return tok[0] + "T" + tok[1] + "Z", err
}

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
    addTimeRange(hostName, TimeRange{st, et})
    return nil
}

func parseCsv(fileName string) error {
    csvFile, err := os.Open(fileName)
    if err != nil {
        fmt.Println(err)
        return err
    }

    reader := csv.NewReader(bufio.NewReader(csvFile))
    for {
        line, error := reader.Read()
        if error == io.EOF {
            break
        } else if error != nil {
            fmt.Println(err)
            return err
        }

        parseParams(line[0], line[1], line[2])
    }
    return nil
}
