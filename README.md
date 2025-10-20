[![Go](https://github.com/msoap/t-plot/actions/workflows/go.yml/badge.svg)](https://github.com/msoap/t-plot/actions/workflows/go.yml)

# t-plot

It's simple utility for plot bar chart from numbers in the terminal.
By default it's plot chart from stdin, one number per line and add bar chart to the end of line.

## Installation

```bash
go install github.com/msoap/t-plot@latest
```
## Options

 - `-k N` - column number for plot, starting from 1, by default try to detect the first column of numbers
 - `-s ...` - style, "bar-simple", "bar-horizontal-1px", "bar-vertical-1px" (default: "bar-simple")
 - `-c "#"` - chart symbol (default: `#`)
 - `-w N` - width of chart (default: rest of terminal width)
 - `-h` - print help and exit

## Examples

```bash
# plot chart from stdin by file size in ls (column 5)
$ ls -l | t-plot -k 5

# chart by file sizes from du
$ du -s some_path/* | t-plot
```

## Screenshot

<img width="824" alt="image" src="https://github.com/msoap/t-plot/assets/844117/3be0f2db-d272-4cc6-a73e-ab5172ba18f5">
