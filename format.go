package strtotime

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	reSpace    = "[ ]+"
	reSpaceOpt = "[ ]*"
	reMeridian = "(am|pm)"
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

func formats() []format {

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
			return r.time(12, 0, 0, 0)
		},
	}

	midnightOrToday := format{
		regex: `(midnight|today)`,
		name:  "midnight | today",
		callback: func(r *result, inputs ...string) error {
			return r.resetTime()
		},
	}

	tomorrow := format{
		regex: "(tomorrow)",
		name:  "tomorrow",
		callback: func(r *result, inputs ...string) error {
			r.rd++
			// Original code calls r.resetTime() here.
			return nil
		},
	}

	timestamp := format{
		regex: `^@(-?\d+)`,
		name:  "timestamp",
		callback: func(r *result, inputs ...string) error {
			s, err := strconv.Atoi(inputs[0])

			if err != nil {
				return err
			}

			r.rs += s
			r.y = pointer(1970)
			r.m = pointer(0)
			r.d = pointer(1)
			r.dates = 0

			return r.resetTime()
			// original code called r.zone(0)
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
		regex: "(?i)" + "^(" + reMonthFull + "|" + reMonthAbbr + ")",
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
		regex: "^" + reHour24 + ":" + reMinutelz + ":" + reSecondlz + "[:.]([0-9]+)" + reMeridian + "?",
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

			return r.time(hour, minute, second, frac)
		},
	}

	timeLong12 := format{
		regex: "^" + reHour12 + "[:.]" + reMinute + "[:.]" + reSecondlz + reSpaceOpt + reMeridian,
		name:  "timeLong12",
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

			meridian := inputs[3]

			return r.time(processMeridian(hour, meridian), minute, second, 0)
		},
	}

	timeShort12 := format{
		regex: "^" + reHour12 + "[:.]" + reMinutelz + reSpaceOpt + reMeridian,
		name:  "timeShort12",
		callback: func(r *result, inputs ...string) error {

			hour, err := strconv.Atoi(inputs[0])
			if err != nil {
				return err
			}

			minute, err := strconv.Atoi(inputs[1])
			if err != nil {
				return err
			}

			meridian := inputs[2]

			return r.time(processMeridian(hour, meridian), minute, 0, 0)
		},
	}

	timeTiny12 := format{
		regex: "^" + reHour12 + reSpaceOpt + reMeridian,
		name:  "timeTiny12",
		callback: func(r *result, inputs ...string) error {

			hour, err := strconv.Atoi(inputs[0])
			if err != nil {
				return err
			}

			meridian := inputs[1]

			return r.time(processMeridian(hour, meridian), 0, 0, 0)
		},
	}

	soap := format{
		regex: "^" + reYear4 + "-" + reMonthlz + "-" + reDaylz + "T" + reHour24lz + ":" + reMinutelz + ":" + reSecondlz + reFrac + "(z|Z)?" + reTzCorrection + "?",
		name:  "soap",
		callback: func(r *result, inputs ...string) error {

			year, err := strconv.Atoi(inputs[0])
			if err != nil {
				return err
			}
			month, err := strconv.Atoi(inputs[1])
			if err != nil {
				return err
			}
			day, err := strconv.Atoi(inputs[2])
			if err != nil {
				return err
			}
			hour, err := strconv.Atoi(inputs[3])
			if err != nil {
				return err
			}

			minute, err := strconv.Atoi(inputs[4])
			if err != nil {
				return err
			}

			second, err := strconv.Atoi(inputs[5])
			if err != nil {
				return err
			}

			mili := inputs[6]

			if len(mili) > 3 {
				mili = mili[0:3]
			}

			frac, err := strconv.Atoi(mili)
			if err != nil {
				return err
			}

			tzCorrection := inputs[8]

			err = r.ymd(year, month-1, day)
			if err != nil {
				return err
			}
			err = r.time(hour, minute, second, frac)
			if err != nil {
				return err
			}
			if len(tzCorrection) > 0 {
				r.zone(processTzCorrection(tzCorrection, 0))
			}
			return nil
		},
	}

	wddx := format{
		regex: "^" + reYear4 + "-" + reMonth + "-" + reDay + "T" + reHour24 + ":" + reMinute + ":" + reSecond,
		name:  "wddx",
		callback: func(r *result, inputs ...string) error {

			year, err := strconv.Atoi(inputs[0])
			if err != nil {
				return err
			}
			month, err := strconv.Atoi(inputs[1])
			if err != nil {
				return err
			}
			day, err := strconv.Atoi(inputs[2])
			if err != nil {
				return err
			}
			hour, err := strconv.Atoi(inputs[3])
			if err != nil {
				return err
			}

			minute, err := strconv.Atoi(inputs[4])
			if err != nil {
				return err
			}

			second, err := strconv.Atoi(inputs[5])
			if err != nil {
				return err
			}

			err = r.ymd(year, month-1, day)
			if err != nil {
				return err
			}

			err = r.time(hour, minute, second, 0)
			return err
		},
	}

	exif := format{
		regex: "(?i)" + "^" + reYear4 + ":" + reMonthlz + ":" + reDaylz + " " + reHour24lz + ":" + reMinutelz + ":" + reSecondlz,
		name:  "exif",
		callback: func(r *result, inputs ...string) error {
			year, err := strconv.Atoi(inputs[0])
			if err != nil {
				return err
			}
			month, err := strconv.Atoi(inputs[1])
			if err != nil {
				return err
			}
			day, err := strconv.Atoi(inputs[2])
			if err != nil {
				return err
			}
			hour, err := strconv.Atoi(inputs[3])
			if err != nil {
				return err
			}

			minute, err := strconv.Atoi(inputs[4])
			if err != nil {
				return err
			}

			second, err := strconv.Atoi(inputs[5])
			if err != nil {
				return err
			}

			err = r.ymd(year, month-1, day)
			if err != nil {
				return err
			}

			err = r.time(hour, minute, second, 0)
			return err
		},
	}

	formats := []format{
		yesterday,
		now,
		noon,
		midnightOrToday,
		tomorrow,
		timestamp,
		firstOrLastDay,
		monthFullOrMonthAbbr,
		mssqltime,
		timeLong12,
		timeShort12,
		timeTiny12,
		soap,
		wddx,
		exif,
	}

	return formats
}
