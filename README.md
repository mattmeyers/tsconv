# tsconv

`tsconv` is a command line utiliy for converting timestamps to different formats.

## Installation

```sh
go install github.com/mattmeyers/tsconv
```

## Usage

```
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
```

## Examples

- Convert an epoch timestamp to an MST RFC822 timestamp

    ```sh
    $ tsconv --out RFC822 --tz MST 1671849943
    23 Dec 22 19:45 MST
    ```

- Convert a UTC timestamp to EST using an offset

    ```sh
    $ tsconv --tz -5 2022-12-24T02:47:52Z
    2022-12-23T21:47:52-05:00
    ```
