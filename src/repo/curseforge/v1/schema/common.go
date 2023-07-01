package schema

import (
	"strings"
	"time"
)

// Reference: https://docs.curseforge.com/#tocS_ModLoaderType
type ModLoaderType int

func (l ModLoaderType) String() string {
	switch int(l) {
	case int(Forge):
		return "Forge"
	case int(Cauldron):
		return "Cauldron"
	case int(LiteLoader):
		return "LiteLoader"
	case int(Fabric):
		return "Fabric"
	case int(Quilt):
		return "Quilt"
	default:
		return "Unknown"
	}
}

func ModLoader(modLoader string) ModLoaderType {
	switch strings.ToLower(modLoader) {
	case "forge":
		return Forge
	case "cauldron":
		return Cauldron
	case "liteloader":
		return LiteLoader
	case "fabric":
		return Fabric
	case "quilt":
		return Quilt
	default:
		return Any
	}
}

const (
	Any        ModLoaderType = 0
	Forge      ModLoaderType = 1
	Cauldron   ModLoaderType = 2
	LiteLoader ModLoaderType = 3
	Fabric     ModLoaderType = 4
	Quilt      ModLoaderType = 5
)

// Reference: https://docs.curseforge.com/#tocS_HashAlgo

type HashAlgo int

const (
	Sha1 HashAlgo = 1
	Md5  HashAlgo = 2
)

// Reference: https://docs.curseforge.com/#tocS_Pagination
type Pagination struct {
	Index       int32 `json:"index"`       // A zero based index of the first item that is included in the response
	PageSize    int32 `json:"pageSize"`    // The requested number of items to be included in the response
	ResultCount int32 `json:"resultCount"` // The actual number of items that were included in the response
	TotalCount  int64 `json:"totalCount"`  // The total number of items available by the request
}

// Reference: https://docs.curseforge.com/#tocS_SortableGameVersion
type SortableGameVersion struct {
	GameVersionName        string    `json:"gameVersionName"`        // Original version name (e.g. 1.5b)
	GameVersionPadded      string    `json:"gameVersionPadded"`      // Used for sorting (e.g. 0000000001.0000000005)
	GameVersion            string    `json:"gameVersion"`            // game version clean name (e.g. 1.5)
	GameVersionReleaseDate time.Time `json:"gameVersionReleaseDate"` // Game version release date
	GameVersionTypeID      int       `json:"gameVersionTypeId"`      // Game version type id
}
