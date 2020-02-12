package logfrog

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	dumpColorLabel   = color.New(color.BgBlack).Add(color.FgHiWhite).Add(color.Bold)
	dumpColorString  = color.New(color.BgBlack).Add(color.FgBlue).Add(color.Bold)
	dumpColorNumber  = color.New(color.BgBlack).Add(color.FgGreen).Add(color.Bold)
	dumpColorBoolean = color.New(color.BgBlack).Add(color.FgYellow).Add(color.Bold)
	dumpColorNull    = color.New(color.BgBlack).Add(color.FgWhite).Add(color.Bold)
)

func dump(left func(line int), v interface{}, indent int, label string, line int) {
	if label != "" {
		left(line)
		dumpColorLabel.Print(strings.Repeat("  ", indent))
		dumpColorLabel.Print(label)
	} else {
		dumpColorLabel.Print(strings.Repeat("  ", indent))
	}
	switch t := v.(type) {
	case string:
		dumpColorString.Println("\"" + strings.ReplaceAll(v.(string), "\n", "\\n") + "\"")
	case float64, float32, int, int64, int32, int16, int8, uint, uint16, uint32, uint64:
		dumpColorNumber.Println(v)
	case bool:
		if v.(bool) {
			dumpColorBoolean.Println("true")
		} else {
			dumpColorBoolean.Println("false")
		}
	case nil:
		fmt.Println("null")
	case LogData:
		for k, value := range v.(LogData) {
			dump(left, value, indent+1, k+": ", line+1)
		}
	case map[string]interface{}:
		if len(v.(map[string]interface{})) > 0 {
			fmt.Println()
		}
		for k, value := range v.(map[string]interface{}) {
			dump(left, value, indent+1, k+": ", line+1)
		}
	case map[string]string:
		if len(v.(map[string]string)) > 0 {
			fmt.Println()
		}
		for k, value := range v.(map[string]string) {
			dump(left, value, indent+1, k+": ", line+1)
		}
	case []interface{}:
		sliceValue := v.([]interface{})
		if len(sliceValue) > 0 {
			fmt.Println()
		}
		for _, value := range sliceValue {
			dump(left, value, indent+1, "- ", line+1)
		}
	default:
		fmt.Printf("I don't know about type %T!\n", t)
	}
}
