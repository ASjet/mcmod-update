package v1

import (
	"encoding/json"
	"fmt"
	"mcmod-update/src/repo/curseforge/v1/schema"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	APIKEY_ENV = "CURSE_FORGE_APIKEY"
)

var (
	cli  = testClient()
	skip = false // Set to true to skip tests
)

func TestGetModFiles(t *testing.T) {
	if skip {
		t.Skipf("env %q is empty, skipped", APIKEY_ENV)
	}

	files, err := cli.GetModFiles(223794, 0, 1)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	printJson(files)
}

func testClient() *Client {
	key := os.Getenv(APIKEY_ENV)
	if len(key) == 0 || skip {
		skip = true
	}

	return NewClient(key, "1.19.2", schema.Forge)
}

func printJson(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
