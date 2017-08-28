package utils

import (
	"fmt"

	"time"
)

const (
	ymdhmszFmt = "2006-01-02 15:04:05-07"
	ymdFmt     = "2006-01-02"
)

func TimeNow() string {
	// UTC: "Etc/GMT"
	// Shanghai: "Asia/Shanghai"
	// For more abbrs, please look into golang source code:
	//     go/src/time/zoneinfo_abbrs_windows.go
	loc, err := time.LoadLocation("Etc/GMT")
	if err != nil {
		//fmt.Printf("util", "unable to get timelocation from string: %s", err)
		return ymdhmszFmt
	}
	return time.Now().In(loc).Format(ymdhmszFmt)
}

func DateStr() string {
	// UTC: "Etc/GMT"
	// Shanghai: "Asia/Shanghai"
	// For more abbrs, please look into golang source code:
	//     go/src/time/zoneinfo_abbrs_windows.go
	loc, err := time.LoadLocation("Etc/GMT")
	if err != nil {
		fmt.Printf("util", "unable to get timelocation from string: %s", err)
	}
	return time.Now().In(loc).Format(ymdFmt)
}

func FormatTime(t time.Time) string {
	loc, _ := time.LoadLocation("Etc/GMT")
	return t.In(loc).Format(ymdhmszFmt)
}
