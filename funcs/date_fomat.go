package funcs

import (
	"strconv"
	"time"
)

func IsLTMDate(d string) bool {
	lastYear := time.Now().AddDate(-1, 0, 0)
	if len(d) == 4 {
		year, err := strconv.Atoi(d)
		if err != nil {
			return false
		}
		if year >= lastYear.Year() {
			return true
		}
	} else if len(d) == 8 {
		timeString := "20060102"
		loc, _ := time.LoadLocation("Local")
		t, err := time.ParseInLocation(timeString, d, loc)
		if err != nil {
			return false
		}
		if t.Unix() >= lastYear.Unix() {
			return true
		}
	}

	return false
}
