package model

import "time"

type File struct {
	ModID        int32     `json:"mod_id"`
	FileID       int32     `json:"file_id"`
	DispName     string    `json:"display_name"`
	FileName     string    `json:"file_name"`
	ReleaseType  string    `json:"release_type"`
	Hash         string    `json:"hash"`
	Date         time.Time `json:"date"`
	DownloadUrl  string    `json:"download_url"`
	GameVersions []string  `json:"game_versions"`
	McVersion    string    `json:"mc_version"`
	ModLoader    string    `json:"mod_loader"`
	RequiredDeps []int32   `json:"required_deps"`
	OptionalDeps []int32   `json:"optional_deps"`
}
