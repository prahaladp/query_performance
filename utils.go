package main

import (
    "time"
    "fmt"
)

func convertTimeToString(st time.Time) string {
    return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
        st.Year(), st.Month(), st.Day(),
        st.Hour(), st.Minute(), st.Second())
}
