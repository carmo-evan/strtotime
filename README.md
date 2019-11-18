[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

# strtotime

![Image of PHP mascot being switch by Golang mascot](https://i.imgur.com/8RhHjkD.jpg)

Golang implementation of `strtotime`, a very popular PHP function for converting English text to a timestamp. This is an exercise inspired on [this](https://github.com/kvz/locutus/blob/master/src/php/datetime/strtotime.js) Javascript implementation.

## Installation

`strtotime` uses Go modules and is compatible with Go 1.12 upwards. Install using `go get`.

```
go get gitub.com/carmo-evan/strtotime
```

After importing it, the `strtotime` package will expose only one method: `Parse`. It takes two arguments - an English string describing some point in time; and a unix timestamp that should represent the current time, or another referencial point in time you want to use. 

Try it on [the playground](https://play.golang.org/p/gfqdPaE6XlU).

```go
package main

import (
	"fmt"
	"github.com/carmo-evan/strtotime"
	"time"
)

func main() {
    //Now is Nov 17, 2019
    u, err := strtotime.Parse("next Friday 3pm", time.Now().Unix())
    
    if err != nil {
        // crash and burn
    }

    t := time.Unix(u,0)
    
    fmt.Println(t)
    //output: 2019-11-22 15:00:00 +0000 UTC
}
```

## Supported Formats

- [x] yesterday
- [x] now
- [x] noon
- [x] midnightOrToday
- [x] tomorrow
- [x] timestamp
- [x] firstOrLastDay
- [ ] backOrFrontOf
- [ ] weekdayOf
- [x] mssqltime
- [x] timeLong12
- [x] timeShort12
- [x] timeTiny12
- [x] soap
- [x] wddx
- [x] exif
- [x] xmlRpc
- [x] xmlRpcNoColon
- [x] clf
- [x] iso8601long
- [x] dateTextual
- [x] pointedDate4
- [x] pointedDate2
- [x] timeLong24
- [x] dateNoColon
- [x] pgydotd
- [x] timeShort24
- [x] iso8601noColon
- [x] dateSlash
- [x] american
- [x] americanShort
- [x] gnuDateShortOrIso8601date2
- [x] iso8601date4
- [x] gnuNoColon
- [x] gnuDateShorter
- [x] pgTextReverse
- [x] dateFull
- [x] dateNoDay
- [x] dateNoDayRev
- [x] pgTextShort
- [x] dateNoYear
- [x] dateNoYearRev
- [x] isoWeekDay
- [x] relativeText
- [x] relative
- [x] dayText
- [x] relativeTextWeek
- [x] monthFullOrMonthAbbr
- [x] tzCorrection
- [x] ago
- [x] gnuNoColon2
- [x] year4
- [x] whitespace

### Author

By [Evan do Carmo](https://github.com/carmo-evan)