package logfrog

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dumpData = `{
	"foo": "bar",
	"slice": [
		0,
		1,
		"two",
		{
			"type": "test"
		}
	],
	"sepp": {
		"a": 1.3,
		"true" : true,
		"false": false,
		"null": null
	}
}`

func TestDump(t *testing.T) {
	data := map[string]interface{}{}
	assert.NoError(t, json.Unmarshal([]byte(dumpData), &data))
	dump(func(line int) {}, data, 0, "sepp", 0)
}
