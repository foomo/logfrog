package logfrog

import (
	"encoding/json"
)

type ReaderGograpple struct{}

func (pr *ReaderGograpple) Valid(line string) bool {
	if line == "" {
		return false
	}
	return []byte(line)[0] == byte('{')
}

func (pr *ReaderGograpple) Read(line string) (label string, logData LogData, err error) {
	return readGograppleLine(line)
}

func readGograppleLine(line string) (label string, logData LogData, err error) {
	logData = LogData{}
	errParse := json.Unmarshal([]byte(line), &logData)
	if errParse != nil {
		err = errParse
		return
	}

	switch tm := logData["msg"].(type) {
	case string:
		label = "dlv"
		if len(tm) > 2 {
			msgBytes := []byte(tm)
			if msgBytes[0] == '{' && msgBytes[len(msgBytes)-1] == '}' {
				label = "process"
				logData = LogData{}
				errParse := json.Unmarshal(msgBytes, &logData)
				if errParse != nil {
					err = errParse
					return
				}
			}
		}
	}

	return
}
