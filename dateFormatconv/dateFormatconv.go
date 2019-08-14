package dateFormatconv

import (
	"github.com/pkg/errors"
	"strings"
)

var (
	stdLongMonth             = []byte("January")
	stdMonth                 = []byte("Jan")
	stdNumMonth              = []byte("1")
	stdZeroMonth             = []byte("01")
	stdLongWeekDay           = []byte("Monday")
	stdWeekDay               = []byte("Mon")
	stdDay                   = []byte("2")
	stdUnderDay              = []byte("_2")
	stdZeroDay               = []byte("02")
	stdHour                  = []byte("15")
	stdHour12                = []byte("3")
	stdZeroHour12            = []byte("03")
	stdMinute                = []byte("4")
	stdZeroMinute            = []byte("04")
	stdSecond                = []byte("5")
	stdZeroSecond            = []byte("05")
	stdLongYear              = []byte("2006")
	stdYear                  = []byte("06")
	stdPM                    = []byte("PM")
	stdpm                    = []byte("pm")
	stdTZ                    = []byte("MST")
	stdISO8601TZ             = []byte("Z0700") // prints Z for UTC
	stdISO8601SecondsTZ      = []byte("Z070000")
	stdISO8601ShortTZ        = []byte("Z07")
	stdISO8601ColonTZ        = []byte("Z07:00") // prints Z for UTC
	stdISO8601ColonSecondsTZ = []byte("Z07:00:00")
	stdNumTZ                 = []byte("-0700") // always numeric
	stdNumSecondsTz          = []byte("-070000")
	stdNumShortTZ            = []byte("-07")    // always numeric
	stdNumColonTZ            = []byte("-07:00") // always numeric
	stdNumColonSecondsTZ     = []byte("-07:00:00")
	stdFracSecond0           = []byte(".0") //, ".00", ... , trailing zeros included
	stdFracSecond9           = []byte(".9") //, ".99", ..., trailing zeros omitted
)

const (
	status_closed = iota
	status_open
)

var (
	runeMap = map[byte]map[int][]byte{
		121: map[int][]byte{
			2: stdYear,
			4: stdLongYear,
		},
		77: map[int][]byte{
			1: stdNumMonth,
			2: stdZeroMonth,
			3: stdMonth,
			4: stdLongMonth,
		},
		109: map[int][]byte{
			1: stdMinute,
			2: stdZeroMinute,
		},
		// 68:  "D",  day of the year
		100: map[int][]byte{
			1: stdDay,
			2: stdZeroDay,
		},
		72: map[int][]byte{
			1: stdHour,
			2: stdHour,
		},
		83: map[int][]byte{
			1: precisionMillisecond(1),
			2: precisionMillisecond(2),
			3: precisionMillisecond(3),
		},
		115: map[int][]byte{
			1: stdSecond,
			2: stdZeroSecond,
		},
		69: map[int][]byte{
			1: stdWeekDay,
			2: stdWeekDay,
			3: stdWeekDay,
			4: stdLongWeekDay,
		},
		// 101: "e",
		90: map[int][]byte{
			1: stdISO8601TZ,
			2: stdISO8601ColonTZ,
		},
		// 122: "z",
	}
)

var (
	ERROR_InvalidFormatString = errors.New("invalid format string.")
)

func precisionMillisecond(l int) []byte {
	return []byte(strings.Repeat("9", l))
}

type pointer struct {
	count  int
	rune   *byte
	status int
}

func (p *pointer) reset() {
	p.count = 0
	p.status = status_closed
	p.rune = nil
}

func (p *pointer) end() ([]byte, error) {
	if p.status == status_closed {
		return nil, nil
	}
	return p.format()
}

func (p *pointer) format() ([]byte, error) {
	defer p.reset()
	fl, ok := runeMap[*p.rune]
	if !ok {
		return nil, ERROR_InvalidFormatString
	}
	bs, ok := fl[p.count]
	if !ok {
		if *p.rune == 83 {
			return precisionMillisecond(p.count), nil
		}
		return nil, ERROR_InvalidFormatString
	}
	return bs, nil
}

func (p *pointer) read(r byte) ([]byte, error) {
	switch p.status {
	case status_open:
		if r != *p.rune {
			bs, err := p.format()
			if err != nil {
				return nil, err
			}
			nbs, err := p.read(r)
			if err != nil {
				return nil, err
			}
			bs = append(bs, nbs...)
			return bs, nil
		} else {
			p.count += 1
		}
	case status_closed:
		if _, ok := runeMap[r]; ok {
			p.status = status_open
			p.rune = &r
			p.count = 1
		} else {
			return []byte{r}, nil
		}
	}
	return nil, nil
}

func Format(s string) (string, error) {
	bs := []byte{}
	p := &pointer{}
	for _, r := range []byte(s) {
		nbs, err := p.read(r)
		if err != nil {
			return "", err
		}
		bs = append(bs, nbs...)
	}
	nbs, err := p.end()
	if err != nil {
		return "", err
	}
	return string(append(bs, nbs...)), nil
}
