package logfrog

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const simpleLine = `{"level":"info"}`

func TestRead(t *testing.T) {
	_, _, errReadCrap := read("[crap")
	require.Error(t, errReadCrap)
	_, logData, errRead := read(simpleLine)
	require.NoError(t, errRead)
	assert.Equal(t, "info", logData["level"])
}

func TestReadDockerComposeLine(t *testing.T) {
	expectedLabel := "some sörvice"
	label, logData, errRead := readDockerComposeLine(" " + expectedLabel + " 	| 	" + simpleLine)
	assert.Equal(t, expectedLabel, label)
	require.NoError(t, errRead)
	assert.Equal(t, "info", logData["level"])
}

func TestReadGoGrappleLine(t *testing.T) {
	expectedMessage := "just a message"
	msgJSON := `{"msg":"` + expectedMessage + `"}`
	label, logData, errRead := readGograppleLine(msgJSON)
	assert.Equal(t, "dlv", label)
	require.NoError(t, errRead)
	assert.Equal(t, expectedMessage, logData["msg"])
	nested := strings.ReplaceAll(msgJSON, `"`, "\\\"")
	processMsg := fmt.Sprintf(`{"msg": "%s"}`, nested)
	label, logData, errRead = readGograppleLine(processMsg)
	assert.Equal(t, "process", label)
	require.NoError(t, errRead)
	assert.Equal(t, expectedMessage, logData["msg"])
}
