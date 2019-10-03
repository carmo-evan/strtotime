package strtotime

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	reSpace    = "[ \\t]+"
	reSpaceOpt = "[ \\t]*"
	reMeridian = "(?:([ap])\\.?m\\.?([\\t ]|$))"
	reHour24   = "(2[0-4]|[01]?[0-9])"
	reHour24lz = "([01][0-9]|2[0-4])"
	reHour12   = "(0?[1-9]|1[0-2])"
	reMinute   = "([0-5]?[0-9])"
	reMinutelz = "([0-5][0-9])"
	reSecond   = "(60|[0-5]?[0-9])"
	reSecondlz = "(60|[0-5][0-9])"
	reFrac     = "(?:\\.([0-9]+))"

	reDayfull = "sunday|monday|tuesday|wednesday|thursday|friday|saturday"
	reDayabbr = "sun|mon|tue|wed|thu|fri|sat"
	reDaytext = reDayfull + "|" + reDayabbr + "|weekdays?"

	reReltextnumber = "first|second|third|fourth|fifth|sixth|seventh|eighth?|ninth|tenth|eleventh|twelfth"
	reReltexttext   = "next|last|previous|this"
	reReltextunit   = "(?:second|sec|minute|min|hour|day|fortnight|forthnight|month|year)s?|weeks|" + reDaytext

	reYear          = "([0-9]{1,4})"
	reYear2         = "([0-9]{2})"
	reYear4         = "([0-9]{4})"
	reYear4withSign = "([+-]?[0-9]{4})"
	reMonth         = "(1[0-2]|0?[0-9])"
	reMonthlz       = "(0[0-9]|1[0-2])"
	reDay           = "(?:(3[01]|[0-2]?[0-9])(?:st|nd|rd|th)?)"
	reDaylz         = "(0[0-9]|[1-2][0-9]|3[01])"

	reMonthFull  = "january|february|march|april|may|june|july|august|september|october|november|december"
	reMonthAbbr  = "jan|feb|mar|apr|may|jun|jul|aug|sept?|oct|nov|dec"
	reMonthroman = "i[vx]|vi{0,3}|xi{0,2}|i{1,3}"
	reMonthText  = "(" + reMonthFull + "|" + reMonthAbbr + "|" + reMonthroman + ")"

	reTzCorrection = "((?:GMT)?([+-])" + reHour24 + ":?" + reMinute + "?)"
	reDayOfYear    = "(00[1-9]|0[1-9][0-9]|[12][0-9][0-9]|3[0-5][0-9]|36[0-6])"
	reWeekOfYear   = "(0[1-9]|[1-4][0-9]|5[0-3])"
)

type format struct {
	regex    string
	name     string
	callback func(r *result, inputs ...string) error
}

func pointer(x int) *int {
	return &x
}

func formats() map[string]format {

	yesterday := format{
		regex: `(yesterday)`,
		name:  "yesterday",
		callback: func(r *result, inputs ...string) error {
			r.rd--
			//HACK: Original code had call to r.resetTime()
			// Might have to do with timezone adjustment
			return nil
		},
	}

	now := format{
		regex: `(now)`,
		name:  "now",
		callback: func(r *result, inputs ...string) error {
			return nil
		},
	}

	noon := format{
		regex: `(noon)`,
		name:  "noon",
		callback: func(r *result, inputs ...string) error {
			r.resetTime()
			r.time(12, 0, 0, 0)
			return nil
		},
	}

	midnightOrToday := format{
		regex: `(midnight|today)`,
		name:  "midnight | today",
		callback: func(r *result, inputs ...string) error {
			r.resetTime()
			return nil
		},
	}

	tomorrow := format{
		regex: "(tomorrow)",
		name:  "tomorrow",
		callback: func(r *result, inputs ...string) error {
			r.rd++
			// r.resetTime()
			return nil
		},
	}

	timestamp := format{
		regex: `^@(-?\d+)`,
		name:  "timestamp",
		callback: func(r *result, inputs ...string) error {
			s, err := strconv.Atoi(inputs[0])
			r.rs += s
			r.y = pointer(1970)
			r.m = pointer(0)
			r.d = pointer(1)
			r.dates = 0
			r.resetTime()
			// r.zone(0)
			return err
		},
	}

	firstOrLastDay := format{
		regex: `^(first|last) day of`,
		name:  "firstdayof | lastdayof",
		callback: func(r *result, inputs ...string) error {
			if strings.ToLower(inputs[0]) == "first" {
				r.firstOrLastDayOfMonth = 1
				return nil
			}
			r.firstOrLastDayOfMonth = -1
			return nil
		},
	}

	monthFullOrMonthAbbr := format{
		regex: "(?i)^(january|february|march|april|may|june|july|august|september|october|november|december|jan|feb|mar|apr|may|jun|jul|aug|sept?|oct|nov|dec)",
		name:  "monthfull | monthabbr",
		callback: func(r *result, inputs ...string) error {
			month := inputs[0]
			if r.dates > 0 {
				return fmt.Errorf("strtotime: The string contains two conflicting date/months")
			}
			r.dates++
			r.m = pointer(lookupMonth(month))
			return nil
		},
	}

	// weekdayOf := format{
	// 	regex: "^(reReltextnumber|reReltexttext)(reDayfull|reDayabbr) of",
	// 	name: "weekdayof",
	// 	callback: func(r *result, inputs ...string) error {

	// 	},
	// 	//TODO:Implement
	//   },

	mssqltime := format{
		regex: "^(2[0-4]|[01]?[0-9]):([0-5][0-9]):(60|[0-5][0-9])[:.]([0-9]+)(am|pm)?",
		name:  "mssqltime",
		callback: func(r *result, inputs ...string) error {

			hour, err := strconv.Atoi(inputs[0])
			if err != nil {
				return err
			}

			minute, err := strconv.Atoi(inputs[1])
			if err != nil {
				return err
			}

			second, err := strconv.Atoi(inputs[2])
			if err != nil {
				return err
			}

			frac, err := strconv.Atoi(inputs[3][0:3])
			if err != nil {
				return err
			}
			if len(inputs) == 5 {
				meridian := inputs[4]
				hour = processMeridian(hour, meridian)
			}
			r.time(hour, minute, second, frac)
			return nil
		},
	}

	formats := map[string]format{
		"yesterday":            yesterday,
		"now":                  now,
		"noon":                 noon,
		"midnightOrToday":      midnightOrToday,
		"tomorrow":             tomorrow,
		"timestamp":            timestamp,
		"firstOrLastDay":       firstOrLastDay,
		"monthFullOrMonthAbbr": monthFullOrMonthAbbr,
		"mssqltime":            mssqltime,
	}

	return formats
}
