package strtotime

import (
	"testing"
	"time"
)

var now = time.Now()

var parseTests = []struct {
	in      string
	out     int64
	success bool
}{
	{"yesterday noon", time.Date(now.Year(), now.Month(), now.Day()-1, 12, 0, 0, 0, time.UTC).Unix(), true},
	{"now", time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC).Unix(), true},
	{"midnight", time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Unix(), true},
	{"tomorrow", time.Date(now.Year(), now.Month(), now.Day()+1, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC).Unix(), true},
	{"@1569600000", 1569600000, true},
	{"last day of October", time.Date(now.Year(), time.October, 31, 0, 0, 0, 0, time.UTC).Unix(), true},
	{"01:59:59.040", time.Date(now.Year(), now.Month(), now.Day(), 1, 59, 59, 40000000, time.UTC).Unix(), true},
	{"01:59:59.040pm", time.Date(now.Year(), now.Month(), now.Day(), 13, 59, 59, 40000000, time.UTC).Unix(), true},
	{"16:59:59.040", time.Date(now.Year(), now.Month(), now.Day(), 16, 59, 59, 40000000, time.UTC).Unix(), true},
	{"01:59:59 pm", time.Date(now.Year(), now.Month(), now.Day(), 13, 59, 59, 0, time.UTC).Unix(), true},
	{"01:59:59am", time.Date(now.Year(), now.Month(), now.Day(), 1, 59, 59, 0, time.UTC).Unix(), true},
	{"01.59.59pm", time.Date(now.Year(), now.Month(), now.Day(), 13, 59, 59, 0, time.UTC).Unix(), true},
	{"01:59 pm", time.Date(now.Year(), now.Month(), now.Day(), 13, 59, 0, 0, time.UTC).Unix(), true},
	{"01:59pm", time.Date(now.Year(), now.Month(), now.Day(), 13, 59, 0, 0, time.UTC).Unix(), true},
	{"01.59pm", time.Date(now.Year(), now.Month(), now.Day(), 13, 59, 0, 0, time.UTC).Unix(), true},
	{"01 pm", time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, time.UTC).Unix(), true},
	{"01am", time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, time.UTC).Unix(), true},
	{"tomorrow 01am", time.Date(now.Year(), now.Month(), now.Day()+1, 1, 0, 0, 0, time.UTC).Unix(), true},
	{"last day of October 1am", time.Date(now.Year(), time.October, 31, 1, 0, 0, 0, time.UTC).Unix(), true},
	{"1am 2pm", 0, false},
	{"2008-10-31T15:07:38.6875000-05:00", 1225483658, true},
	{"2008-10-31T15:07:38.034567890GMT-05:00", 1225483658, true},
	{"2008-10-31T15:07:38.034567890Z", 1225465658, true},
	{"2008-10-31T15:07:38", 1225465658, true},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		t.Run(tt.in, func(t *testing.T) {
			r, err := Parse(tt.in)
			if err != nil && tt.success {
				t.Fatal(err)
			}
			if r != tt.out && tt.success {
				t.Errorf("Result should have been %v, but it was %v", tt.out, r)
			}
		})
	}
}

func TestProcessMeridian(t *testing.T) {
	h := processMeridian(12, "am")
	if h != 0 {
		t.Errorf("h should've been 0, but it is %v", h)
	}
	h = processMeridian(10, "pm")
	if h != 22 {
		t.Errorf("h should've been 22, but it is %v", h)
	}
}

var monthTests = []struct {
	in  string
	out int
}{
	{"jan", 0},
	{"january", 0},
	{"i", 0},
	{"feb", 1},
	{"february", 1},
	{"ii", 1},
	{"mar", 2},
	{"march", 2},
	{"iii", 2},
	{"apr", 3},
	{"april", 3},
	{"iv", 3},
	{"may", 4},
	{"v", 4},
	{"jun", 5},
	{"june", 5},
	{"vi", 5},
	{"jul", 6},
	{"july", 6},
	{"vii", 6},
	{"aug", 7},
	{"august", 7},
	{"viii", 7},
	{"sep", 8},
	{"sept", 8},
	{"september", 8},
	{"ix", 8},
	{"oct", 9},
	{"october", 9},
	{"x", 9},
	{"nov", 10},
	{"november", 10},
	{"xi", 10},
	{"dec", 11},
	{"december", 11},
	{"xii", 11},
}

func TestLookupMonth(t *testing.T) {
	for _, tt := range monthTests {
		t.Run(tt.in, func(t *testing.T) {
			m := lookupMonth(tt.in)
			if m != tt.out {
				t.Errorf("Output should've been %v, but it was %v.", tt.out, m)
			}
		})
	}
}

var yearTests = []struct {
	in  string
	out int
}{
	{"2020", 2020},
	{"75", 1975},
	{"20", 2020},
	{"2002", 2002},
}

func TestProcessYear(t *testing.T) {
	for _, tt := range yearTests {
		t.Run(tt.in, func(t *testing.T) {
			y, err := processYear(tt.in)
			if err != nil {
				t.Error(err)
			}
			if y != tt.out {
				t.Errorf("Output should've been %v, but it was %v.", tt.out, y)
			}
		})
	}
}

var weekdayTests = []struct {
	in  string
	out int
}{
	{"mon", 1},
	{"monday", 1},
	{"tue", 2},
	{"tuesday", 2},
	{"wed", 3},
	{"wednesday", 3},
	{"thu", 4},
	{"thursday", 4},
	{"fri", 5},
	{"friday", 5},
	{"sat", 6},
	{"saturday", 6},
	{"sun", 0},
	{"sunday", 0},
}

func TestLookupWeekday(t *testing.T) {
	for _, tt := range weekdayTests {
		t.Run(tt.in, func(t *testing.T) {
			d := lookupWeekday(tt.in, 0)
			if d != tt.out {
				t.Errorf("Output should've been %v, but it was %v.", tt.out, d)
			}
		})
	}
}

var relativeTests = []struct {
	in  string
	out map[string]int
}{
	{"last", map[string]int{"amount": -1, "behavior": 0}},
	{"previous", map[string]int{"amount": -1, "behavior": 0}},
	{"this", map[string]int{"amount": 0, "behavior": 1}},
	{"first", map[string]int{"amount": 1, "behavior": 0}},
	{"next", map[string]int{"amount": 1, "behavior": 0}},
	{"second", map[string]int{"amount": 2, "behavior": 0}},
	{"third", map[string]int{"amount": 3, "behavior": 0}},
	{"fourth", map[string]int{"amount": 4, "behavior": 0}},
	{"fifth", map[string]int{"amount": 5, "behavior": 0}},
	{"sixth", map[string]int{"amount": 6, "behavior": 0}},
	{"seventh", map[string]int{"amount": 7, "behavior": 0}},
	{"eight", map[string]int{"amount": 8, "behavior": 0}},
	{"eighth", map[string]int{"amount": 8, "behavior": 0}},
	{"ninth", map[string]int{"amount": 9, "behavior": 0}},
	{"tenth", map[string]int{"amount": 10, "behavior": 0}},
	{"eleventh", map[string]int{"amount": 11, "behavior": 0}},
	{"twelfth", map[string]int{"amount": 12, "behavior": 0}},
}

func TestLookupRelative(t *testing.T) {
	for _, tt := range relativeTests {
		t.Run(tt.in, func(t *testing.T) {
			d := lookupRelative(tt.in)
			if d["amount"] != tt.out["amount"] || d["behavior"] != tt.out["behavior"] {
				t.Errorf("Output should've been %v, but it was %v.", tt.out, d)
			}
		})
	}
}

var tzCorrectionTests = []struct {
	in  string
	out int
}{
	{"GMT-5", 300},
	{"GMT-5:00", 300},
	{"GMT-5:30", 330},
	{"GMT+5", -300},
	{"GMT+5:00", -300},
	{"GMT+5:30", -330},
}

func TestTzCorrection(t *testing.T) {
	for _, tt := range tzCorrectionTests {
		t.Run(tt.in, func(t *testing.T) {
			y := processTzCorrection(tt.in, 0)
			if y != tt.out {
				t.Errorf("Output should've been %v, but it was %v.", tt.out, y)
			}
		})
	}
}

var resultToDateTests = []struct {
	n   string
	r   result
	out int64
}{
	{"Sep 27 2019, 8:30am", result{
		y: pointer(2019),
		h: pointer(8),
		m: pointer(8),
		i: pointer(30),
		d: pointer(27),
		s: pointer(0),
	}, 1569573000},
	{"March 1st 2020, 8:30am", result{
		y: pointer(2020),
		h: pointer(8),
		m: pointer(2),
		i: pointer(30),
		d: pointer(1),
		s: pointer(0),
	}, 1583051400},
}

func TestResultToDate(t *testing.T) {
	for _, tt := range resultToDateTests {
		t.Run(tt.n, func(t *testing.T) {
			u := tt.r.toDate().Unix()
			if u != tt.out {
				t.Errorf("Unix stamp should've been %v but it was %v", tt.out, u)
			}
		})
	}
}
