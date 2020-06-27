package logfrog

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/nsf/termbox-go"
)

type Printer struct {
	labelWidth        int
	colorsForServices map[string]color.Attribute
	lastLabel         string
	timestamps        bool
	w                 int
}

const maxLabelWidth = 40

var allColors = []color.Attribute{
	color.BgBlue,
	color.BgCyan,
	color.BgGreen,
	color.BgHiBlue,
	color.BgHiCyan,
	color.BgHiGreen,
	color.BgHiMagenta,
	color.BgHiRed,
	color.BgHiYellow,
}

func getWidth() int {
	// this is a bit clunky, but maybe it is x-platform ...
	errInit := termbox.Init()
	if errInit != nil {
		return 80
	}
	w, _ := termbox.Size()
	termbox.Close()
	return w
}

// NewCLIPrinter
func NewCLIPrinter(timestamps bool) (*Printer, error) {
	p := &Printer{timestamps: timestamps, labelWidth: 20, colorsForServices: map[string]color.Attribute{}, w: getWidth()}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c)
		for {
			switch <-c {
			case syscall.SIGWINCH:
				p.w = getWidth()
				halfWidth := p.w / 2
				if p.labelWidth > halfWidth {
					p.labelWidth = halfWidth
				}
			}
		}
	}()
	return p, nil
}

func (p *Printer) Print(label string, logData LogData) {
	p.block(label, p.lastLabel, logData)
	p.lastLabel = label
}

func (p *Printer) colorForService(name string) color.Attribute {
	col, colOk := p.colorsForServices[name]
	if colOk {
		return col
	}
	// gather some stats
	stats := map[color.Attribute]int{}
	for _, colorForService := range allColors {
		stats[colorForService] = 0
	}
	for _, colorForService := range p.colorsForServices {
		stats[colorForService]++
	}
	minUse := 999999999999999999
	for _, use := range stats {
		if use < minUse {
			minUse = use
		}
	}
	for _, colorForService := range allColors {
		if stats[colorForService] == minUse {
			p.colorsForServices[name] = colorForService
			break
		}
	}
	return p.colorsForServices[name]
}

func (p *Printer) Error(arg ...interface{}) {
	color.New(color.BgBlack).Add(color.FgHiRed).Println(arg...)
}
func (p *Printer) Info(arg ...interface{}) {
	color.New(color.BgBlack).Add(color.FgHiWhite).Println(arg...)
}

func (p *Printer) block(label string, lastLabel string, logData LogData) {
	// some colors
	colorNormal := color.New(color.BgBlack).Add(color.FgWhite)
	//colorDump := color.New(color.BgBlack).Add(color.FgHiWhite).Add(color.Bold)
	colorForServiceLine := color.New(p.colorForService(label)).Add(color.FgWhite)
	colorForService := colorForServiceLine.Add(color.Bold)

	// is there enough place for my label
	labelWidth := 0
	for range label {
		labelWidth++
	}
	if labelWidth > p.labelWidth {
		halfWidth := p.w / 2
		if labelWidth < halfWidth {
			p.labelWidth = labelWidth
		} else if p.labelWidth < halfWidth {
			p.labelWidth = halfWidth
		}
	}

	//dataBlock := ""
	trimmedLabel := " "
	trimmedLabelLength := 1
	const padding = 2
	if label != "" {
		for i, r := range label {
			if i == p.labelWidth {
				break
			}
			trimmedLabel += string(r)
			trimmedLabelLength++
		}
		for trimmedLabelLength < p.labelWidth+padding {
			trimmedLabelLength++
			trimmedLabel += " "
		}
	}

	const colSep = " "
	rightWidth := p.w - (p.labelWidth + padding + 3)

	leftBlock := ""
	leftSep := ""
	if label != "" {
		leftBlock = strings.Repeat(" ", p.labelWidth+padding)
		leftSep = strings.Repeat("-", p.labelWidth+padding)
	}

	// extract some data
	logLevel := strings.ToLower(extract(logData, "level", "Level", "logLevel"))
	logTime := extract(logData, "time", "timestamp", "Timestamp")
	logMsg := strings.Trim(extract(logData, "msg", "Message", "message"), "\n")
	logStack := extract(logData, "stack")

	// level color
	colorLevel := color.New(color.BgBlack).Add(color.Bold)
	switch logLevel {
	case "info":
		colorLevel.Add(color.FgHiGreen)
	case "warning", "warn":
		colorLevel.Add(color.FgHiYellow)
	case "fatal", "error", "exception":
		colorLevel.Add(color.FgHiRed)
	default:
		colorLevel.Add(color.FgHiWhite)
	}

	left := func(line int) {
		if line == 0 {
			if lastLabel != label {
				colorForService.Print(trimmedLabel)
			} else {
				colorForService.Print(leftSep)
			}
		} else {
			colorForService.Print(leftBlock)
		}
	}
	if logLevel != "" {
		// we got a log level
		block(left, func(line int, linePart string) {
			colorLevel.Println(linePart)
		}, " "+logLevel+": "+logMsg, rightWidth)
	} else {
		// not a properly readable one line log msg
		block(left, func(line int, linePart string) {
			colorNormal.Print(colSep)
			colorLevel.Println(linePart)
		}, logMsg, rightWidth)
	}

	// log time if set
	if logTime == "" && p.timestamps {
		logTime = fmt.Sprint(time.Now())
	}
	if logTime != "" {
		colorForService.Print(leftBlock)
		colorNormal.Print(colSep)
		colorNormal.Println(logTime)
	}

	// anything else to log
	if len(logData) > 0 {
		dump(func(line int) {
			colorForService.Print(leftBlock)
		}, logData, 0, "", 0)
	}

	// and a stock at the very end
	if len(logStack) > 0 {
		multiBlock(func(line int) {
			if label != "" {
				colorForService.Print(leftBlock)
				colorNormal.Print(colSep)
			}
		}, func(line int, linePart string) {
			colorNormal.Println(linePart)
		}, logStack, rightWidth)

	}
}
func multiBlock(left func(int), right func(int, string), lines string, width int) {
	for _, line := range strings.Split(lines, "\n") {
		block(left, right, line, width)
	}
}
func block(left func(int), right func(int, string), line string, width int) {
	linePart := ""
	i := 0
	lineI := 0
	for _, r := range line {
		linePart += string(r)
		i++
		if i > width {
			if linePart != "" {
				left(lineI)
				right(lineI, linePart)
			}
			lineI++
			i = 0
			linePart = ""
		}
	}
	if linePart != "" {
		left(lineI)
		right(lineI, linePart)
	}
}

func extract(logData LogData, fields ...string) string {
	for _, field := range fields {
		d, ok := logData[field]
		if ok {
			delete(logData, field)
			return fmt.Sprint(d)
		}
	}
	return ""
}
