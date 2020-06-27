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
	msg, okMSGOk := logData["msg"]
	if okMSGOk {
		switch msg.(type) {
		case string:
			label = "dlv"
			msgString := msg.(string)
			if len(msgString) > 2 {
				msgBytes := []byte(msgString)
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
	}
	return
}
