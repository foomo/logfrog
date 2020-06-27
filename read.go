package logfrog

import (
	"encoding/json"
)

type ReaderType string

const (
	ReaderTypePlain         ReaderType = ""
	ReaderTypeDockerCompose ReaderType = "docker-compose"
	ReaderTypeStern         ReaderType = "stern"
	ReaderTypeGograpple     ReaderType = "gograpple"
)

func GetAvailableTypes() []ReaderType {
	return []ReaderType{ReaderTypeDockerCompose, ReaderTypeStern, ReaderTypeGograpple}
}

type LogData map[string]interface{}

type LogReader interface {
	Valid(line string) bool
	Read(line string) (label string, logData LogData, err error)
}

type ReaderPlain struct{}

func (pr *ReaderPlain) Valid(line string) bool {
	return line != ""
}

func (pr *ReaderPlain) Read(line string) (label string, logData LogData, err error) {
	return read(line)
}

// Read read a json log line
func read(line string) (label string, logData LogData, errRead error) {
	label = ""
	errParse := json.Unmarshal([]byte(line), &logData)
	if errParse != nil {
		errRead = errParse
		return
	}
	return
}
