package strtotime

import (
	"math"
	"time"
)

//result holds all the integers tha make up the final Time object returned.
// we use pointers for some properties because we need to verify if they'be been
// initialized or not
type result struct {
	// date
	y *int
	m *int
	d *int
	// time
	h *int
	i *int
	s *int

	// relative shifts
	ry int
	rm int
	rd int
	rh int
	ri int
	rs int
	rf int

	// weekday related shifts
	weekday         *int
	weekdayBehavior int

	// first or last day of month
	// 0 none, 1 first, -1 last
	firstOrLastDayOfMonth int

	// timezone correction in minutes
	// z *int

	// counters
	dates int
	times int
	zones int
}

func (r *result) ymd(y, m, d int) bool {
	if r.dates > 0 {
		return false
	}

	r.dates++
	*r.y = y
	*r.m = m
	*r.d = d
	return true
}

func (r *result) time(h, i, s, f int) bool {
	if r.times > 0 {
		return false
	}

	r.times++
	*r.h = h
	*r.i = i
	*r.s = s

	return true
}

func (r *result) resetTime() bool {
	*r.h = 0
	*r.i = 0
	*r.s = 0
	r.times = 0

	return true
}

// func (r *result) zone(minutes int) bool {
// 	if r.zones <= 1 {
// 		r.zones++
// 		*r.z = minutes
// 		return true
// 	}
// 	return false
// }

func (r *result) toDate() time.Time {

	relativeTo := time.Now()
	if r.dates > 0 && r.times <= 0 {
		*r.h = 0
		*r.i = 0
		*r.s = 0
	}

	// fill holes
	if r.y == nil {
		y := relativeTo.Year()
		r.y = &y
	}

	if r.m == nil {
		m := lookupMonth(relativeTo.Month().String())
		r.m = &m
	}

	if r.d == nil {
		d := relativeTo.Day()
		r.d = &d
	}

	if r.h == nil {
		h := relativeTo.Hour()
		r.h = &h
	}

	if r.i == nil {
		i := relativeTo.Minute()
		r.i = &i
	}

	if r.s == nil {
		s := relativeTo.Second()
		r.s = &s
	}

	// adjust special early
	switch r.firstOrLastDayOfMonth {
	case 1:
		*r.d = 1
		break
	case -1:
		*r.d = 0
		*r.m++
		break
	}

	if r.weekday != nil {

		var dow = lookupWeekday(relativeTo.Weekday().String(), 1)

		if r.weekdayBehavior == 2 {
			// To make "r week" work, where the current day of week is a "sunday"
			if dow == 0 && *r.weekday != 0 {
				*r.weekday = -6
			}

			// To make "sunday r week" work, where the current day of week is not a "sunday"
			if *r.weekday == 0 && dow != 0 {
				*r.weekday = 7
			}

			*r.d -= dow
			*r.d += *r.weekday
		} else {
			var diff = *r.weekday - dow

			// some PHP magic
			if (r.rd < 0 && diff < 0) || (r.rd >= 0 && diff <= -r.weekdayBehavior) {
				diff += 7
			}

			if *r.weekday >= 0 {
				*r.d += diff
			} else {
				//TODO: Fix this madness
				*r.d -= int((7 - (math.Abs(float64(*r.weekday)) - float64(dow))))
			}

			r.weekday = nil
		}
	}

	// adjust relative
	*r.y += r.ry
	*r.m += r.rm
	*r.d += r.rd

	*r.h += r.rh
	*r.i += r.ri
	*r.s += r.rs

	r.ry = 0
	r.rm = 0
	r.rd = 0
	r.rh = 0
	r.ri = 0
	r.rs = 0
	r.rf = 0

	// note: this is done twice in PHP
	// early when processing special relatives
	// and late
	// todo: check if the logic can be reduced
	// to just one time action
	switch r.firstOrLastDayOfMonth {
	case 1:
		*r.d = 1
		break
	case -1:
		m := lookupNumberToMonth(*r.m)
		firstOfMonth := time.Date(*r.y, m, 1, 0, 0, 0, 0, time.UTC)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		*r.m = lookupMonth(lastOfMonth.String())
		break
	}

	// TODO: process and adjust timezone
	// if r.z != nil {
	// 	*r.i += *r.z
	// }

	return time.Date(*r.y, lookupNumberToMonth(*r.m), *r.d, *r.h, *r.i, *r.s, 0, time.UTC)
}
