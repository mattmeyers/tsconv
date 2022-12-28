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

## Formatting

On top of the predefined standards listed under the supported formats, `tsconv` provides a custom formatting language via the `tsfmt` package. This language is similar to that of PHP and JS, except for some changes that enforce consistency. A `tsfmt` format string is a series of characters that correspond to pieces of a timestamp. The examples in the following table all correspond to the RFC3339 timestamp `2022-01-02T16:01:02-04:00`.

| Symbol | Description                    | Example   |
| ------ | ------------------------------ | --------- |
| `Y`    | Four character year            | `2022`    |
| `y`    | Two character year             | `22`      |
| `M`    | Short month name               | `Jan`     |
| `MM`   | Long month name                | `January` |
| `m`    | Month index                    | `1`       |
| `mm`   | Zero padded month number       | `01`      |
| `D`    | Short day name                 | `Sun`     |
| `DD`   | Long day name                  | `Sunday`  |
| `d`    | Day number                     | `2`       |
| `dd`   | Zero padded day number         | `02`      |
| `H`    | 24-hour clock hour             | `16`      |
| `HH`   | Zero padded 24-hour clock hour | `04`      |
| `h`    | 12-hour clock hour             | `4`       |
| `hh`   | 12-hour clock hour             | `04`      |
| `i`    | Minutes                        | `1`       |
| `ii`   | Zero padded minutes            | `01`      |
| `s`    | Seconds                        | `2`       |
| `ss`   | Zero padded seconds            | `02`      |
| `z`    | Timezone offset                | `-04:00`  |
| `Z`    | Timezone name                  | `EST`     |

Any remaining characters not in this table will be treated as literals and appear in the final string. If a character is reserved, prepending it with a backslash will print the character as a literal. For example, the RFC3339 format can be written as `Y-mm-ddTHH:ii:ssz`.

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
