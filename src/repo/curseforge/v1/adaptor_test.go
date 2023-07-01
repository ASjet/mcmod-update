package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var adp = NewAdaptor(apiKey())

func TestGetModLatestFile(t *testing.T) {
	if skip {
		t.Skipf("env %q is empty, skipped", APIKEY_ENV)
	}

	file, err := adp.GetLatestModFile(223794, "1.19.2", "forge")
	assert.NoError(t, err)
	printJson(file)
}

func TestGetModLatestFileWithDeps(t *testing.T) {
	if skip {
		t.Skipf("env %q is empty, skipped", APIKEY_ENV)
	}

	files, err := adp.GetLatestModFileWithDeps(501214, "1.19.2", "forge", false)
	assert.NoError(t, err)
	assert.Len(t, files, 3)

	files, err = adp.GetLatestModFileWithDeps(501214, "1.19.2", "forge", true)
	assert.NoError(t, err)
	assert.Len(t, files, 5)
	printJson(files)
}
