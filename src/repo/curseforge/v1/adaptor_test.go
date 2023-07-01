package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var adp = NewAdaptor(testClient())

func TestGetModLatestFile(t *testing.T) {
	if skip {
		t.Skipf("env %q is empty, skipped", APIKEY_ENV)
	}

	file, err := adp.GetModLatestFile(223794)
	assert.NoError(t, err)
	printJson(file)
}

func TestGetModLatestFileWithDeps(t *testing.T) {
	if skip {
		t.Skipf("env %q is empty, skipped", APIKEY_ENV)
	}

	files, err := adp.GetModLatestFileWithDeps(501214, false)
	assert.NoError(t, err)
	assert.Len(t, files, 3)

	files, err = adp.GetModLatestFileWithDeps(501214, true)
	assert.NoError(t, err)
	assert.Len(t, files, 5)
	printJson(files)
}
