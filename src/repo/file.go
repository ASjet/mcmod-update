package repo

import "mcmod-update/src/model"

type FileFetcher interface {
	Version() string
	ModLoader() string
	GetLatestModFile(modId int32) (*model.File, error)
	GetLatestModFileWithDeps(modId int32, optional bool) ([]*model.File, error)
}
