package repo

import "mcmod-update/src/model"

type FileFetcher interface {
	GetLatestModFile(modId int32, version, modLoader string) (*model.File, error)
	GetLatestModFileWithDeps(modId int32, version, modLoader string,
		optional bool) ([]*model.File, error)
}
