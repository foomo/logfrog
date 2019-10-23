package logfrog

import (
	"encoding/json"
	"errors"
	"strings"
)

// ReadDockerComposeLine reads a line from a docker-compose log output
func ReadDockerComposeLine(line string) (
	label string,
	logData LogData,
	errRead error,
) {
	found := false
	trimmedLine := ""
	label = ""
	for _, r := range line {
		if found {
			trimmedLine += string(r)
		} else {
			if r == '|' {
				found = true
			} else {
				label += string(r)
			}
		}
	}
	if !found {
		errRead = errors.New("| not found")
	}
	if strings.HasPrefix(trimmedLine, " ") {
		// bit hacky ...
		trimmedLine = strings.TrimPrefix(trimmedLine, " ")
	}
	label = strings.Trim(label, "	 ")
	logData, errRead = Read(strings.Trim(trimmedLine, " 	\n"))
	if errRead != nil {
		logData = map[string]interface{}{"msg": trimmedLine}
		errRead = nil
	}
	return
}

// Read read a json log line
func Read(line string) (logData map[string]interface{}, errRead error) {
	errParse := json.Unmarshal([]byte(line), &logData)
	if errParse != nil {
		errRead = errParse
		return
	}
	return
}
