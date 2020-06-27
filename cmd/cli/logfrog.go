package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/foomo/logfrog"
)

var version string = "dev"

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

const usage = `
usage examples:

# docker-compose:
docker-compose logs -f --tail 1 --no-color | logfrog -log-type docker-compose

# stern:
stern -o json -n some-name-space | logfrog -log-type stern

# gograpple:
gograpple delve my-deployment -c my-container -v --input main.go -n my-ns --vscode | logfrog -log-type gograpple
`

func main() {

	cli := flag.NewFlagSet("logfrog", flag.ExitOnError)
	flagJS := cli.String("js-filter", "", "/path/to/file with js function filter")
	flagHelp := cli.Bool("help", false, "show help")
	flagVersion := cli.Bool("version", false, "show version")
	flagTimestamps := cli.Bool("timestamps", false, "add timestamps, if not present")
	flagLogType := cli.String("log-type", string(logfrog.ReaderTypePlain), "docker-compose | stern | gograpple")

	errParse := cli.Parse(os.Args[1:])

	if errParse != nil || *flagHelp {
		fmt.Println(os.Args[0])
		cli.PrintDefaults()
		fmt.Println(usage)
		os.Exit(1)
	}

	if *flagVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	logType := logfrog.ReaderType(*flagLogType)

	var logReader logfrog.LogReader
	switch logType {
	case logfrog.ReaderTypeDockerCompose:
		logReader = &logfrog.ReaderDockerCompose{}
	case logfrog.ReaderTypeStern:
		logReader = &logfrog.ReaderStern{}
	case logfrog.ReaderTypePlain:
		logReader = &logfrog.ReaderPlain{}
	case logfrog.ReaderTypeGograpple:
		logReader = &logfrog.ReaderGograpple{}
	default:
		must("unknown log type: '"+string(logType)+"'", errors.New("log type must be one of: "+strings.Join(func() (types []string) {
			for _, t := range logfrog.GetAvailableTypes() {
				types = append(types, string(t))
			}
			return types
		}(), ", ")))
	}

	filter, filterModTime, errWatchJSFilter := watchJSFilter(*flagJS, time.Time{})
	must("could not load filter", errWatchJSFilter)

	reader := bufio.NewReader(os.Stdin)
	printer, errPrinter := logfrog.NewCLIPrinter(*flagTimestamps)

	must("could not initialize printer", errPrinter)

	for {
		line, errReadString := reader.ReadString('\n')
		if errReadString != nil {
			fmt.Println("done reading", errReadString)
			os.Exit(1)
		}
		// trim stuff
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, " ")
		line = strings.Trim(line, "	")

		// skip if not valid
		if !logReader.Valid(line) {
			continue
		}

		// run the configured reader
		label, logData, errRead := logReader.Read(line)
		if errRead != nil {
			logData = logfrog.LogData{"msg": line}
		}

		// js filter update?
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

		// execute filter
		errFilter := filter(label, &logData)
		if errFilter != nil {
			fmt.Println("filter error", errFilter)
		}

		// print filtered result
		if logData != nil && len(logData) > 0 {
			printer.Print(label, logData)
		}
	}
}
