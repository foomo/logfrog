package logfrog

import (
	"encoding/json"
)

type ReaderType string

const (
	ReaderTypePlain         ReaderType = ""
	ReaderTypeDockerCompose ReaderType = "docker-compose"
	ReaderTypeStern         ReaderType = "stern"
)

func GetAvailableTypes() []ReaderType {
	return []ReaderType{ReaderTypeDockerCompose, ReaderTypeStern}
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
func read(line string) (_ string, logData LogData, errRead error) {
	_ = ""
	errParse := json.Unmarshal([]byte(line), &logData)
	if errParse != nil {
		errRead = errParse
		return
	}
	return
}
