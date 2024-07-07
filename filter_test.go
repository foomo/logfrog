package logfrog

import (
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GetCurrentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
func TestFilter(t *testing.T) {
	filterFunc, errGetFilter := GetFilter(path.Join(GetCurrentDir(), "test", "example-filter-level-info.js"))
	require.NoError(t, errGetFilter)
	ld := LogData{"level": "error"}
	require.NoError(t, filterFunc("foo", &ld))
	assert.Len(t, ld, 1)
	ldInfo := LogData{"level": "info"}
	require.NoError(t, filterFunc("foo", &ldInfo))
	assert.Nil(t, ldInfo)
}

func TestFilterBrokenJS(t *testing.T) {
	_, errGetFilter := GetFilter(path.Join(GetCurrentDir(), "test", "broken-js-filter.js"))
	assert.Error(t, errGetFilter)
}
