package repo

import "strings"

const (
	ModLoaderForge  = "forge"
	ModLoaderFabric = "fabric"
)

type Version struct {
	McVersion string
	ModLoader string
}

func NewVersion(mcVersion, modLoader string) *Version {
	return &Version{
		McVersion: strings.ToLower(mcVersion),
		ModLoader: strings.ToLower(modLoader),
	}
}
