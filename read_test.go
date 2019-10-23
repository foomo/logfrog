package logfrog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const simpleLine = `{"level":"info"}`

func TestRead(t *testing.T) {
	_, errReadCrap := Read("[crap")
	assert.Error(t, errReadCrap)
	logData, errRead := Read(simpleLine)
	assert.NoError(t, errRead)
	assert.Equal(t, "info", logData["level"])
}

func TestReadDockerComposeLine(t *testing.T) {
	expectedLabel := "some sörvice"
	label, logData, errRead := ReadDockerComposeLine(" " + expectedLabel + " 	| 	" + simpleLine)
	assert.Equal(t, expectedLabel, label)
	assert.NoError(t, errRead)
	assert.Equal(t, "info", logData["level"])
}
