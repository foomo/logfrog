package logfrog

import (
	"strings"
)

type ReaderDockerCompose struct {
}

func (pr *ReaderDockerCompose) Valid(line string) bool {
	if line == "" {
		return false
	}
	parts := strings.Split(line, "|")
	return len(parts) != 1
}

func (pr *ReaderDockerCompose) Read(line string) (label string, logData LogData, err error) {
	return readDockerComposeLine(line)
}

// ReadDockerComposeLine reads a line from a docker-compose log output
func readDockerComposeLine(line string) (
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
		logData = map[string]interface{}{"msg": "| not found"}
		return
	}

	// bit hacky ...
	trimmedLine = strings.TrimPrefix(trimmedLine, " ")

	label = strings.Trim(label, "	 ")
	_, logData, errRead = read(strings.Trim(trimmedLine, " 	\n"))
	if errRead != nil {
		logData = map[string]interface{}{"msg": trimmedLine}
		errRead = nil
	}
	return
}
