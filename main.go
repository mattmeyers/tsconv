package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	errNoInput       = errors.New("no input provided")
	errUnableToParse = errors.New("unable to parse input")
)

var timeFormats = []string{
	time.RFC3339,
	time.RFC822,
	time.UnixDate,
}

const helpMessage = `
Usage:  tsconv [--out OUT] [TIMESTAMP]

Convert a timestamp to another format

The provided timestamp can be provided in any of the supported formats listed
below. Inputs can be supplied as the argument or via stdin through a pipe or
file redirect. If no timestamp is provided, then the current time will be used.

Options:
  --out     The output format to use. Check below for allowed values.
            (default: rfc3339)
  --tz      The output timezone to use. This can be any standard IANA timezone
            or a +/- offset. (default: UTC)

Supported Formats:
  rfc3339   2006-01-02T10:04:05-05:00
  rfc822    02 Jan 06 15:04 MST
  unix      Mon Jan  2 10:04:05 EST 2006
  epoch     1136214245

Examples:
  - Convert an epoch timestamp to an MST RFC822 timestamp
      $ tsconv --out RFC822 --tz MST 1671849943
      23 Dec 22 19:45 MST

  - Convert a UTC timestamp to EST using an offset
      $ tsconv --tz -5 2022-12-24T02:47:52Z
      2022-12-23T21:47:52-05:00
`

func main() {
	if err := initializeApp().run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		flag.CommandLine.Usage()
		os.Exit(1)
	}
}

type app struct {
	opts options
	args []string

	clock clock

	r io.Reader
	w io.Writer
}

type clock interface {
	Now() time.Time
}

type stdClock struct{}

func (stdClock) Now() time.Time { return time.Now() }

func (a app) run() error {
	input, err := a.getInput()
	if err != nil && err != errNoInput {
		return err
	}

	inputTime, err := a.parseInput(input)
	if err != nil {
		return err
	}

	inputTime, err = setTimezone(inputTime, a.opts.timezone)
	if err != nil {
		return err
	}

	fmt.Fprintln(a.w, formatOutput(inputTime, a.opts.outputFormat))

	return nil
}

type options struct {
	outputFormat string
	timezone     string
}

func initializeApp() app {
	var opts options

	flag.StringVar(&opts.outputFormat, "out", "rfc3339", "The output format")
	flag.StringVar(&opts.timezone, "tz", "UTC", "The output timezone")

	flag.CommandLine.Usage = func() { fmt.Print(helpMessage) }

	flag.Parse()

	return app{
		opts:  opts,
		clock: stdClock{},
		args:  flag.Args(),
		r:     os.Stdin,
		w:     os.Stdout,
	}

}

func (a app) getInput() (string, error) {
	if len(a.args) > 0 {
		return a.args[0], nil
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if (stat.Mode()&os.ModeNamedPipe) == 0 && !stat.Mode().IsRegular() {
		return "", errNoInput
	}

	// A timestamp should never be more than 256 bytes
	input, err := io.ReadAll(io.LimitReader(a.r, 256))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(input)), nil
}

func (a app) parseInput(s string) (time.Time, error) {
	if s == "" {
		return a.clock.Now(), nil
	}

	out, err := strconv.Atoi(s)
	if err == nil {
		return time.Unix(int64(out), 0), nil
	}

	for _, format := range timeFormats {
		t, err := time.Parse(format, s)
		if err != nil {
			continue
		}

		return t, nil
	}

	return time.Time{}, errUnableToParse
}

func setTimezone(t time.Time, tz string) (time.Time, error) {
	if strings.HasPrefix(tz, "+") || strings.HasPrefix(tz, "-") {
		loc, err := parseOffset(tz)
		if err != nil {
			return time.Time{}, err
		}

		return t.In(loc), nil
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}

	return t.In(loc), nil
}

func parseOffset(tz string) (*time.Location, error) {
	offset, err := strconv.Atoi(tz[1:])
	if err != nil {
		return nil, err
	}

	if tz[0] == '-' {
		offset *= -1
	}

	return time.FixedZone("", int((time.Duration(offset) * time.Hour).Seconds())), nil

}

func formatOutput(t time.Time, format string) string {
	switch strings.ToLower(format) {
	case "rfc822", "822":
		return t.Format(time.RFC822)
	case "rfc3339", "3339":
		return t.Format(time.RFC3339)
	case "unix":
		return t.Format(time.UnixDate)
	case "epoch":
		return strconv.Itoa(int(t.Unix()))
	default:
		return t.Format(time.RFC3339)
	}
}
