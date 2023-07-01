package v1

import (
	"encoding/json"
	"fmt"
	"mcmod-update/src/repo/curseforge/v1/schema"
	"net/http"
)

const (
	BaseUrl = "https://api.curseforge.com"
)

type Client struct {
	modLoader   schema.ModLoaderType
	gameVersion string
	apiKey      string
}

func NewClient(apiKey string, gameVersion string,
	modLoader schema.ModLoaderType) *Client {
	return &Client{
		apiKey:      apiKey,
		modLoader:   modLoader,
		gameVersion: gameVersion,
	}
}

// Reference: https://docs.curseforge.com/#get-mod-file
func (c *Client) GetModFiles(modId int32, pn int, ps int) ([]*schema.File, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/v1/mods/%d/files?gameVersion=%s&modLoaderType=%d&index=%d&pageSize=%d",
			BaseUrl, modId, c.gameVersion, c.modLoader, pn, ps), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", c.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	rsp := new(schema.GetModFilesResponse)
	if err := decoder.Decode(rsp); err != nil {
		return nil, err
	}

	return rsp.Data, nil
}
