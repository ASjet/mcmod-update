package repo

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	BaseUrl    = "https://api.curseforge.com"
	LatestFile = "/v1/mods/%d/files?gameVersion=%s&modLoaderType=%s&index=0&pageSize=1"
)

type Response struct {
	Data []ModFile `json:"data"`
}

type ModFile struct {
	ModId   int    `json:"modId"`
	ModName string `json:"modName"`

	FileId       int       `json:"id,omitempty"`
	DispName     string    `json:"displayName,omitempty"`
	FileName     string    `json:"fileName,omitempty"`
	Date         time.Time `json:"fileDate,omitempty"`
	Url          string    `json:"downloadUrl,omitempty"`
	GameVersions []string  `json:"-"`
	McVersion    string    `json:"mcVersion,omitempty"`
	ModLoader    string    `json:"modLoader,omitempty"`
}

func (m *ModFile) String() string {
	s, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(s)
}

type CurseforgeRepo struct {
	url    string
	apikey string
}

func NewCurseforgeRepo(apikey string) *CurseforgeRepo {
	return &CurseforgeRepo{
		url:    BaseUrl + LatestFile,
		apikey: apikey,
	}
}

func (cf *CurseforgeRepo) LatestModFile(modId int, ver *Version) (*ModFile, error) {
	body, err := NewRequest("GET", fmt.Sprintf(cf.url, modId, ver.McVersion, ver.ModLoader)).
		WithHeader("Accept", "application/json").
		WithHeader("x-api-key", cf.apikey).Do()
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, nil
	}
	rsp := new(Response)
	if err = json.Unmarshal(body, rsp); err != nil {
		return nil, err
	}
	if len(rsp.Data) == 0 {
		return nil, nil
	}
	rsp.Data[0].McVersion = ver.McVersion
	rsp.Data[0].ModLoader = ver.ModLoader
	return &rsp.Data[0], nil
}

type ModFileSlice []*ModFile

func (m ModFileSlice) Len() int {
	return len(m)
}

func (m ModFileSlice) Less(i, j int) bool {
	return m[i].ModId < m[j].ModId
}

func (m ModFileSlice) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
