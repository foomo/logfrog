package logfrog

import (
	"encoding/json"
	"strings"
)

type sternEntry struct {
	Message       string `json:"message"`
	Namespace     string `json:"namespace"`
	PodName       string `json:"podName"`
	ContainerName string `json:"containerName"`
}

type ReaderStern struct{}

func (pr *ReaderStern) Valid(line string) bool {
	return line != "" && line[0] == '{' && line[len(line)-1] == '}'
}

func (pr *ReaderStern) Read(line string) (label string, logData LogData, err error) {
	sd := &sternEntry{}
	logData = LogData{}
	errUnmarshal := json.Unmarshal([]byte(line), &sd)
	if errUnmarshal != nil {
		logData["msg"] = line
		return "unknown", logData, nil //nolint:nilerr
	}
	label = sd.Namespace + ":" + sd.PodName + "(" + sd.ContainerName + ")"
	errLogData := json.Unmarshal([]byte(sd.Message), &logData)
	if errLogData != nil {
		logData["msg"] = strings.Trim(sd.Message, "\n")
	}
	sd.Message = ""
	logData["stern"] = map[string]string{
		"namespace":     sd.Namespace,
		"containerName": sd.ContainerName,
		"podName":       sd.PodName,
	}
	return
}
