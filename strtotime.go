package strtotime

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Parse(s string, relativeTo time.Time) (int64, error) {
	r := &result{}
	formats := formats()
	for {
		noMatch := true
		for _, format := range formats {

			re := regexp.MustCompile(format.regex)
			match := re.FindStringSubmatch(s)

			if len(match) > 0 {

				noMatch = false

				err := format.callback(r, match[1:]...)

				if err != nil {
					return 0, err
				}

				s = strings.TrimSpace(re.ReplaceAllString(s, ""))
				break
			}
		}

		if len(s) == 0 {
			return r.toDate(relativeTo).Unix(), nil
		}

		if noMatch {
			return 0, fmt.Errorf(`strtotime: Unrecognizable input - "%v"`, s)
		}
	}
}

//processMeridian converts 12 hour format type to 24 hour format
func processMeridian(h int, m string) int {
	m = strings.ToLower(m)
	switch m {
	case "am":
		if h == 12 {
			h -= 12
		}
		break
	case "pm":
		if h != 12 {
			h += 12
		}
		break
	}

	return h
}

//processYear converts a year string such as "75" to a year, such as 1975
func processYear(yearStr string) (int, error) {
	y, err := strconv.Atoi(yearStr)

	cutoffYear := 70 //Magic number. Anything before this will be in the 2000s. After, 1900s.

	if err != nil {
		return 0, err
	}

	if len(yearStr) >= 4 || y >= 100 {
		return y, nil
	}

	if y < cutoffYear {
		y += 2000
		return y, nil
	}

	if y >= cutoffYear {
		y += 1900
		return y, nil
	}

	return y, nil
}

func lookupMonth(m string) int {
	monthMap := map[string]int{
		"jan":       0,
		"january":   0,
		"i":         0,
		"feb":       1,
		"february":  1,
		"ii":        1,
		"mar":       2,
		"march":     2,
		"iii":       2,
		"apr":       3,
		"april":     3,
		"iv":        3,
		"may":       4,
		"v":         4,
		"jun":       5,
		"june":      5,
		"vi":        5,
		"jul":       6,
		"july":      6,
		"vii":       6,
		"aug":       7,
		"august":    7,
		"viii":      7,
		"sep":       8,
		"sept":      8,
		"september": 8,
		"ix":        8,
		"oct":       9,
		"october":   9,
		"x":         9,
		"nov":       10,
		"november":  10,
		"xi":        10,
		"dec":       11,
		"december":  11,
		"xii":       11,
	}
	return monthMap[strings.ToLower(m)]
}

func lookupNumberToMonth(m int) time.Month {
	monthMap := map[int]time.Month{
		0:  time.January,
		1:  time.February,
		2:  time.March,
		3:  time.April,
		4:  time.May,
		5:  time.June,
		6:  time.July,
		7:  time.August,
		8:  time.September,
		9:  time.October,
		10: time.November,
		11: time.December,
	}
	return monthMap[m]
}

func lookupWeekday(day string, desiredSundayNumber int) int {
	dayNumberMap := map[string]int{
		"mon":       1,
		"monday":    1,
		"tue":       2,
		"tuesday":   2,
		"wed":       3,
		"wednesday": 3,
		"thu":       4,
		"thursday":  4,
		"fri":       5,
		"friday":    5,
		"sat":       6,
		"saturday":  6,
		"sun":       0,
		"sunday":    0,
	}

	if n, ok := dayNumberMap[strings.ToLower(day)]; ok {
		return n
	}

	return desiredSundayNumber
}

func lookupRelative(rel string) map[string]int {
	relativeNumbersMap := map[string]int{
		"last":     -1,
		"previous": -1,
		"this":     0,
		"first":    1,
		"next":     1,
		"second":   2,
		"third":    3,
		"fourth":   4,
		"fifth":    5,
		"sixth":    6,
		"seventh":  7,
		"eight":    8,
		"eighth":   8,
		"ninth":    9,
		"tenth":    10,
		"eleventh": 11,
		"twelfth":  12,
	}

	relativeBehaviorMap := map[string]int{
		"this": 1,
	}

	relativeBehaviorValue := 0

	if value, ok := relativeBehaviorMap[rel]; ok {
		relativeBehaviorValue = value
	}

	rel = strings.ToLower(rel)

	return map[string]int{
		"amount":   relativeNumbersMap[rel],
		"behavior": relativeBehaviorValue,
	}
}

//processTzCorrection converts a time zone offset (i.e. GMT-5) to minutes (i.e. 300)
func processTzCorrection(tzOffset string, oldValue int) int {
	const reTzCorrectionLoose = `(?:GMT)?([+-])(\d+)(:?)(\d{0,2})`
	re := regexp.MustCompile(reTzCorrectionLoose)
	offsetGroups := re.FindStringSubmatch(tzOffset)

	sign := -1

	if offsetGroups[1] == "-" {
		sign = 1
	}

	hours, err := strconv.Atoi(offsetGroups[2])

	if err != nil {
		return oldValue
	}

	var minutes int

	if len(offsetGroups) > 3 && len(offsetGroups[4]) > 0 {
		minutes, err = strconv.Atoi(offsetGroups[4])

		if err != nil {
			return oldValue
		}
	}

	if len(offsetGroups) < 3 {
		m := float64(hours % 100)
		h := float64(hours / 100)
		minutes = int(math.Floor(m))
		hours = int(math.Floor(h))
	}

	return sign * (hours*60 + minutes)
}
