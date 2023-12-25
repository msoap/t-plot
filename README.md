# t-plot

It's simple utility for plot bar chart from numbers in the terminal.
By default it's plot chart from stdin, one number per line and add bar chart to the end of line.

## Installation

```bash
go install github.com/msoap/t-plot@latest
```

## Usage

```bash
# plot chart from stdin by file size in ls (column 5)
$ ls -l | t-plot -k 5

# chart by file sizes from du
$ du -s some_path/* | t-plot
```

## Options

 - `-k N` - column number for plot (default: 1)
 - `-c "#"` - chart character (default: `#`)
 - `-w N` - width of chart (default: rest of terminal width using $COLUMNS)
 - `-h` - print help and exit
