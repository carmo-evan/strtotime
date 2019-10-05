package strtotime

import (
	"fmt"
	"testing"
)

var formatsMap = formats()

var formatTests = []struct {
	n     string
	f     format
	r     *result
	input string
	out   int64
}{
	{"yesterday", formatsMap["yesterday"], &result{
		y: pointer(2019),
		h: pointer(8),
		m: pointer(8),
		i: pointer(30),
		d: pointer(27),
		s: pointer(0),
	}, "", 1569486600},
	{"noon", formatsMap["noon"], &result{
		y: pointer(2019),
		h: pointer(8),
		m: pointer(8),
		i: pointer(30),
		d: pointer(27),
		s: pointer(0),
	}, "", 1569585600},
	// {"timestamp", formatsMap["timestamp"], &result{
	// 	y: pointer(2019),
	// 	h: pointer(8),
	// 	m: pointer(8),
	// 	i: pointer(30),
	// 	d: pointer(27),
	// 	s: pointer(0),
	}, "@1569600000", 1569600000},
	{"first day", formatsMap["firstOrLastDay"], &result{
		y: pointer(2019),
		h: pointer(8),
		m: pointer(8),
		i: pointer(30),
		s: pointer(0),
	}, "first", 1567326600},
	{"last day", formatsMap["firstOrLastDay"], &result{
		y: pointer(2019),
		h: pointer(8),
		m: pointer(1),
		i: pointer(30),
		s: pointer(0),
	}, "last", 1551342600},
	{"american", formatsMap["american"], &result{
		y: pointer(2019),
		m: pointer(10),
		d: pointer(4),
	}, "10/4/19", 1570147200},
}

func TestFormats(t *testing.T) {
	for _, tt := range formatTests {
		t.Run(tt.n, func(t *testing.T) {
			fmt.Println(tt.r.toDate().Unix())
			tt.f.callback(tt.r, tt.input) //1569475800
			fmt.Println(tt.r.toDate().Unix())
			if tt.r.toDate().Unix() != tt.out {
				t.Errorf("Unix stamp should've been %v but it was %v", tt.out, tt.r.toDate().Unix())
			}
		})
	}
}
