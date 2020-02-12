package logfrog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const simpleLine = `{"level":"info"}`

func TestRead(t *testing.T) {
	_, _, errReadCrap := read("[crap")
	assert.Error(t, errReadCrap)
	_, logData, errRead := read(simpleLine)
	assert.NoError(t, errRead)
	assert.Equal(t, "info", logData["level"])
}

func TestReadDockerComposeLine(t *testing.T) {
	expectedLabel := "some sörvice"
	label, logData, errRead := readDockerComposeLine(" " + expectedLabel + " 	| 	" + simpleLine)
	assert.Equal(t, expectedLabel, label)
	assert.NoError(t, errRead)
	assert.Equal(t, "info", logData["level"])
}
