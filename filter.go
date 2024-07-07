package logfrog

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/robertkrimen/otto"
)

// Filter func for logEntries
type Filter func(service string, ld *LogData) error

// GetFilter get a filter function with a js file, it has to contain a function filter(service, logData) { }
func GetFilter(file string) (filterFunc Filter, err error) {
	fileBytes, errRead := os.ReadFile(file)
	if errRead != nil {
		err = errRead
		return
	}
	vm := otto.New()
	_, errRunFile := vm.Run(string(fileBytes))
	if errRunFile != nil {
		err = errRunFile
		return
	}

	filterFunc = func(service string, ld *LogData) error {
		jsonLogData, errMarshal := json.Marshal(ld)
		if errMarshal != nil {
			return errMarshal
		}
		errSetLogData := vm.Set("logData", string(jsonLogData))
		if errSetLogData != nil {
			return errors.New("could not set logData in filter vm:" + errSetLogData.Error())
		}
		errSetService := vm.Set("service", service)
		if errSetService != nil {
			return errors.New("could not set service in filter vm:" + errSetService.Error())
		}
		res, errRun := vm.Run(`JSON.stringify(filter(JSON.parse(logData), service));`)
		if errRun != nil {
			return errRun
		}
		if res.IsString() {
			filteredLD := LogData{}
			errUnmarshal := json.Unmarshal([]byte(res.String()), &filteredLD)
			if errUnmarshal != nil {
				return errors.New("could not unmarshal result: " + errUnmarshal.Error())
			}
			*ld = filteredLD
			return nil
		}
		return errors.New("missing string in filter vm:" + res.String())
	}
	return
}
