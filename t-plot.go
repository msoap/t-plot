/*
Options:
  - `-k N`   - column number for plot (default: 1)
  - `-c "#"` - chart character (default: `#`)
  - `-w N`   - width of chart (default: rest of terminal width using $COLUMNS)
  - `-h`     - print help and exit
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/msoap/byline"
	"golang.org/x/term"
)

const (
	defaultTermWidth = 80
	maxTermWidth     = 130
)

type opt struct {
	columnN int
	barChar string
	width   int
}

type lineData struct {
	num   int
	width int
}

func main() {
	cfg := parseArgs()

	lines, err := readStdin()
	if err != nil {
		printErr("read stdin: %s\n", err)
	}

	info := getTextInfo(cfg, lines)
	maxs := getAllMax(info)
	chartLines := renderChart(cfg, lines, info, maxs)
	fmt.Println(strings.Join(chartLines, "\n"))
}

func printErr(frmt string, err error) {
	fmt.Fprintf(os.Stderr, frmt, err)
	os.Exit(1)
}

func parseArgs() opt {
	res := opt{}

	doHelp := flag.Bool("h", false, "print help and exit")
	flag.IntVar(&res.columnN, "k", 1, "column number for plot")
	flag.StringVar(&res.barChar, "c", "#", "bar chart character")
	flag.IntVar(&res.width, "w", 0, "width of chart")
	flag.Parse()

	if *doHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	return res
}

func readStdin() ([]string, error) {
	return byline.
		NewReader(os.Stdin).
		MapString(func(in string) string {
			// strip trailing newline
			return strings.TrimRight(in, "\n")
		}).
		ReadAllSliceString()
}

func getTextInfo(cfg opt, lines []string) []lineData {
	res := make([]lineData, len(lines))
	for i, line := range lines {
		res[i].width = utf8.RuneCountInString(line) // TODO: use graphic symbol length

		fields := strings.Fields(line)
		if cfg.columnN > len(fields) {
			continue
		}

		res[i].num, _ = strconv.Atoi(fields[cfg.columnN-1])
	}
	return res
}

func getAllMax(info []lineData) lineData {
	maxNum, maxWidth := 0, 0
	for _, item := range info {
		if item.num > maxNum {
			maxNum = item.num
		}
		if item.width > maxWidth {
			maxWidth = item.width
		}
	}
	return lineData{maxNum, maxWidth}
}

func renderChart(cfg opt, lines []string, info []lineData, maxs lineData) []string {
	termWidth := getTermWidth()

	res := make([]string, len(lines))
	for i, line := range lines {
		if info[i].num == 0 {
			res[i] = line
			continue
		}

		chartWidth := int(float64(info[i].num) / float64(maxs.num) * float64(termWidth-info[i].width))
		if chartWidth == 0 {
			res[i] = line
			continue
		}

		chart := strings.Repeat(cfg.barChar, chartWidth)
		if l := utf8.RuneCountInString(line); l < maxs.width {
			line += strings.Repeat(" ", maxs.width-l)
		}
		res[i] = line + "\t" + chart
	}
	return res
}

func getTermWidth() int {
	// "tput cols"/"stty size"/$COLUMNS is not working in programs
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))

	if width == 0 {
		width = defaultTermWidth
	}
	if width > maxTermWidth {
		width = maxTermWidth
	}

	return width
}
