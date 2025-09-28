/*
Options:
  - `-k N`   - column number for plot (default: 1)
  - `-s ...` - style, "bar-simple", "bar-horizontal-1px", "bar-vertical-1px" (default: "bar-simple")
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
	"github.com/msoap/tcg"
	"github.com/msoap/tcg/turtle"
	"golang.org/x/term"
)

const (
	defaultTermWidth = 80
	maxTermWidth     = 120
	minChartWidth    = 10
	widthReserve     = 8
)

type opt struct {
	style   chartStyle
	columnN int
	barChar string
	width   int
}

type lineData struct {
	num   float64
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
	chartLines := createChart(cfg, lines, info, maxs)
	fmt.Println(strings.Join(chartLines, "\n"))
}

func printErr(frmt string, args ...any) {
	fmt.Fprintf(os.Stderr, frmt, args...)
	os.Exit(1)
}

func parseArgs() opt {
	res := opt{}

	doHelp := flag.Bool("h", false, "print help and exit")
	flag.Var(&res.style, "s", `style, "bar-simple", "bar-horizontal-1px", "bar-vertical-1px" (default: "bar-simple")`)
	flag.IntVar(&res.columnN, "k", 1, "column number for plot")
	flag.StringVar(&res.barChar, "c", "■", "bar chart character")
	flag.IntVar(&res.width, "w", 0, "width of chart")
	flag.Parse()

	if *doHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if len(res.barChar) == 0 && res.style == csBarSimple {
		printErr("bar chart character is empty\n")
	}

	return res
}

func readStdin() ([]string, error) {
	return byline.
		NewReader(os.Stdin).
		MapString(func(in string) string {
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

		res[i].num, _ = strconv.ParseFloat(fields[cfg.columnN-1], 64)
	}

	return res
}

func getAllMax(info []lineData) lineData {
	maxNum, maxWidth := 0.0, 0
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

func createChart(cfg opt, lines []string, info []lineData, maxs lineData) []string {
	termWidth := getTermWidth()

	chartWidth := 0
	if cfg.width > 0 {
		chartWidth = cfg.width
	} else {
		chartWidth = termWidth - maxs.width - widthReserve
		if chartWidth < minChartWidth {
			chartWidth = minChartWidth
		}
	}

	switch cfg.style {
	case csBarSimple:
		barChart := renderChartSimple(cfg.barChar, chartWidth, info, maxs)
		lines = alignTextLines(lines, maxs)

		res := make([]string, len(lines))
		for i := range lines {
			res[i] = lines[i] + "\t" + barChart[i]
		}

		return res
	case csBarHorizontal1px:
		return renderChartBarHorizontal1px(chartWidth, info, maxs)
	case csBarVertical1px:
		return renderChartVertical1px(info, maxs)
	default:
		printErr("style %v is not implemented yet\n", cfg.style)
		return nil
	}
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

func alignTextLines(lines []string, maxs lineData) []string {
	res := make([]string, len(lines))
	for i, line := range lines {
		if l := utf8.RuneCountInString(line); l < maxs.width {
			line += strings.Repeat(" ", maxs.width-l)
		}
		res[i] = line
	}

	return res
}

func renderChartSimple(barChar string, width int, info []lineData, maxs lineData) []string {
	canvas := tcg.NewBuffer(width, len(info))

	for i, item := range info {
		chartWidth := int(float64(item.num) / float64(maxs.num) * float64(width))
		canvas.HLine(0, i, chartWidth, tcg.Black)
	}

	firstRune, _ := utf8.DecodeRuneInString(barChar)
	mode, err := tcg.NewPixelMode(1, 1, []rune{' ', firstRune})
	if err != nil {
		printErr("create pixel mode for %q: %s\n", barChar, err)
	}

	res := canvas.RenderAsStrings(*mode)
	if len(info) != len(res) {
		printErr("something went wrong, len(info) != len(res), %d != %d\n", len(info), len(res))
	}

	return res
}

func renderChartBarHorizontal1px(width int, info []lineData, maxs lineData) []string {
	tcgMode := tcg.Mode2x3
	canvas := tcg.NewBuffer(width*tcgMode.Width(), len(info))

	for i, item := range info {
		barLen := int(float64(item.num) / float64(maxs.num) * float64(width*tcgMode.Width()))
		canvas.HLine(0, i, barLen, tcg.Black)
	}

	res := canvas.RenderAsStrings(tcgMode)

	return res
}

func renderChartVertical1px(info []lineData, maxs lineData) []string {
	const heightInChars = 10
	tcgMode := tcg.Mode2x3
	heightInPx := tcgMode.Height() * heightInChars
	canvas := tcg.NewBuffer(len(info), heightInPx)

	trtl := turtle.New(&canvas)
	for i, item := range info {
		barLen := int(float64(item.num) / float64(maxs.num) * float64(heightInPx))
		trtl.GoToAbs(i, heightInPx).Up(barLen)
	}

	res := canvas.RenderAsStrings(tcgMode)

	return res
}
