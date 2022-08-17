package typegen

import (
	"fmt"
	"time"
)

const (
	TimestampTZFMT = "2006-01-02 15:04:05.000000 +07:00"
	TimeTZFMT      = "15:04:05.000000 +07:00"
	defaultStart   = ""
)

func GenTimestampTZ(nullable bool, ps ...interface{}) string {
	var (
		now = time.Now()

		stYear             = 1970
		stMonth time.Month = 1
		stDay              = 1
		stHour             = 0
		stMin              = 0
		stSec              = 0
		stNsec             = 0
		loc                = now.Location()

		endYear             = now.Year()
		endMonth time.Month = 12
		endDay              = 30
		endHour             = 23
		endMin              = 59
		endSec              = 59
		endNsec             = 999999000
	)

	if len(ps) == 2 {
		st, err := time.Parse(TimestampTZFMT, ps[0].(string))
		if err == nil {
			stYear, stMonth, stDay = st.Date()
			stHour, stMin, stSec, stNsec, loc = st.Hour(), st.Minute(), st.Second(), st.Nanosecond(), st.Location()
		}
		end, err := time.Parse(TimestampTZFMT, ps[1].(string))
		if err == nil {
			endYear, endMonth, endDay = end.Date()
			endHour, endMin, endSec, endNsec = end.Hour(), end.Minute(), end.Second(), end.Nanosecond()
		}
	}

	return nullWrap(nullable, fmt.Sprintf("'%s'", time.Date(
		randIntInRange(stYear, endYear),
		time.Month(randIntInRange(int(stMonth), int(endMonth))),
		randIntInRange(stDay, endDay),
		randIntInRange(stHour, endHour),
		randIntInRange(stMin, endMin),
		randIntInRange(stSec, endSec),
		// 999999000 nanosec == 999.999 milisec
		randIntInRange(stNsec, endNsec),
		loc,
	).Format(TimestampTZFMT)))
}

func GenTimeTZ(nullable bool, ps ...interface{}) string {
	var (
		stHour = 0
		stMin  = 0
		stSec  = 0
		stNsec = 0
		loc    = time.Now().Location()

		endHour = 23
		endMin  = 59
		endSec  = 59
		endNsec = 999999000
	)
	if len(ps) == 2 {
		st, err := time.Parse(TimeTZFMT, ps[0].(string))
		if err == nil {
			stHour, stMin, stSec, stNsec, loc = st.Hour(), st.Minute(), st.Second(), st.Nanosecond(), st.Location()
		}
		end, err := time.Parse(TimeTZFMT, ps[1].(string))
		if err == nil {
			endHour, endMin, endSec, endNsec = end.Hour(), end.Minute(), end.Second(), end.Nanosecond()
		}
	}

	return nullWrap(nullable, fmt.Sprintf("'%s'", time.Date(
		1970,
		1,
		1,
		randIntInRange(stHour, endHour),
		randIntInRange(stMin, endMin),
		randIntInRange(stSec, endSec),
		// 999999000 nanosec == 999.999 milisec
		randIntInRange(stNsec, endNsec),
		loc,
	).Format(TimeTZFMT)))
}
