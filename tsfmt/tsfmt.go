package tsfmt

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var shortDayNames = map[string]string{
	"Sunday":    "Sun",
	"Monday":    "Mon",
	"Tuesday":   "Tue",
	"Wednesday": "Wed",
	"Thursday":  "Thu",
	"Friday":    "Fri",
	"Saturday":  "Sat",
}

var shortMonthNames = [...]string{
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
}

type Format struct {
	raw        string
	fmtString  string
	formatters []timeFormatter
}

func Parse(s string) Format {
	fstring, formatters := parseFormatString(s)
	return Format{
		raw:        s,
		fmtString:  fstring,
		formatters: formatters,
	}
}

func parseFormatString(s string) (string, []timeFormatter) {
	var format strings.Builder
	var formatters []timeFormatter

	var i int
	for i < len(s) {
		inc := 1
		f := "%s"
		var fn timeFormatter

		switch s[i] {
		case 'Y':
			fn = longYear
		case 'y':
			fn = shortYear
		case 'M':
			if peek(s, i) == 'M' {
				fn = longAlphaMonth
				inc++
			} else {
				fn = shortAlphaMonth
			}
		case 'm':
			fn = numericMonth
			if peek(s, i) == 'm' {
				f = "%02s"
				inc++
			}
		case 'D':
			if peek(s, i) == 'D' {
				fn = longAlphaDay
				inc++
			} else {
				fn = shortAlphaDay
			}
		case 'd':
			fn = numericDay
			if peek(s, i) == 'd' {
				f = "%02s"
				inc++
			}
		case 'H':
			fn = longHour
			if peek(s, i) == 'H' {
				f = "%02s"
				inc++
			}
		case 'h':
			fn = shortHour
			if peek(s, i) == 'h' {
				f = "%02s"
				inc++
			}
		case 'i':
			fn = minutes
			if peek(s, i) == 'i' {
				f = "%02s"
				inc++
			}
		case 's':
			fn = seconds
			if peek(s, i) == 's' {
				f = "%02s"
				inc++
			}
		case 'z':
			fn = tzOffset
		case 'Z':
			fn = tzName
		case '\\':
			f = string(s[i+1])
			inc++
		default:
			f = string(s[i])
		}

		i += inc
		format.WriteString(f)
		if fn != nil {
			formatters = append(formatters, fn)
		}
	}

	return format.String(), formatters
}

func (f Format) Format(t time.Time) string {
	if f.fmtString == "" {
		return strconv.Itoa(int(t.Unix()))
	}

	var vals []any
	for _, tf := range f.formatters {
		vals = append(vals, tf(t))
	}

	return fmt.Sprintf(f.fmtString, vals...)
}

type timeFormatter func(time.Time) string

func longYear(t time.Time) string {
	return strconv.Itoa(t.Year())
}

func shortYear(t time.Time) string {
	return lastN(strconv.Itoa(t.Year()), 2)
}

func longAlphaMonth(t time.Time) string {
	return t.Month().String()
}

func shortAlphaMonth(t time.Time) string {
	return shortMonthNames[t.Month()-1]
}

func numericMonth(t time.Time) string {
	return strconv.Itoa(int(t.Month()))
}

func longAlphaDay(t time.Time) string {
	return t.Weekday().String()
}

func shortAlphaDay(t time.Time) string {
	return shortDayNames[t.Weekday().String()]
}

func numericDay(t time.Time) string {
	return strconv.Itoa(t.Day())
}

func longHour(t time.Time) string {
	return strconv.Itoa(t.Hour())
}

func shortHour(t time.Time) string {
	return strconv.Itoa(t.Hour() % 12)
}

func minutes(t time.Time) string {
	return strconv.Itoa(t.Minute())
}

func seconds(t time.Time) string {
	return strconv.Itoa(t.Second())
}

func tzOffset(t time.Time) string {
	_, offset := t.Zone()

	var negative bool
	if offset < 0 {
		negative = true
		offset *= -1
	}

	if negative {
		return fmt.Sprintf("-%02d:00", offset/3600)
	} else {
		return fmt.Sprintf("+%02d:00", offset/3600)
	}
}

func tzName(t time.Time) string {
	name, _ := t.Zone()
	return name
}

func peek(s string, i int) byte {
	if i+1 < len(s) {
		return s[i+1]
	}

	return 0
}

func lastN(s string, n int) string {
	return s[len(s)-n:]
}
