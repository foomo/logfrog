package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/foomo/logfrog"
)

func must(comment string, err error) {
	if err != nil {
		fmt.Println(comment, err)
		os.Exit(1)
	}
}

func watchJSFilter(file string, lastMod time.Time) (filter logfrog.Filter, newLastMod time.Time, err error) {
	filter = func(service string, ld *logfrog.LogData) error {
		return nil
	}
	if file == "" {
		return
	}
	filter = nil
	info, errOpen := os.Stat(file)
	if errOpen != nil {
		err = errOpen
		return
	}
	newLastMod = info.ModTime()
	if newLastMod.UnixNano() > lastMod.UnixNano() {
		f, errFilter := logfrog.GetFilter(file)
		if errFilter != nil {
			err = errFilter
			return
		}
		filter = f
	}
	return
}

func main() {

	flagJS := flag.String("js-filter", "", "/path/to/file with js function filter")
	flagDockerCompose := flag.Bool("docker-compose", false, "docker-compose mode - do not forget to add a --no-color on the docker-compose logs ;)")
	flag.Parse()

	filter, filterModTime, errWatchJSFilter := watchJSFilter(*flagJS, time.Time{})
	must("could not load filter", errWatchJSFilter)

	reader := bufio.NewReader(os.Stdin)
	printer, errPrinter := logfrog.NewCLIPrinter()

	must("could not initialize printer", errPrinter)

	for {
		line, errReadString := reader.ReadString('\n')
		if errReadString != nil {
			fmt.Println("done reading", errReadString)
			os.Exit(1)
		}
		if *flagJS != "" {
			nextFilter, nextFilterModTime, errNextFilter := watchJSFilter(*flagJS, filterModTime)
			if errNextFilter != nil {
				printer.Error("could not load filter", errNextFilter)
			} else if nextFilter != nil {
				filter = nextFilter
				filterModTime = nextFilterModTime
				printer.Info("reloaded filter", filterModTime)
			}
		}
		line = strings.Trim(line, "\n")
		label := ""
		var logData logfrog.LogData
		var errRead error
		if *flagDockerCompose {
			parts := strings.Split(line, "|")
			if len(parts) == 1 {
				continue
			}
			label, logData, errRead = logfrog.ReadDockerComposeLine(line)
		} else {
			logData, errRead = logfrog.Read(line)
			if errRead != nil {
				logData = logfrog.LogData{"msg": line}
			}
		}

		errFilter := filter(label, &logData)
		if errFilter != nil {
			fmt.Println("filter error", errFilter)
		}
		if logData != nil && len(logData) > 0 {
			printer.Print(label, logData)
		}
	}
}
