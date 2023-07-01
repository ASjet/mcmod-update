package v1

import (
	"mcmod-update/src/model"
	"mcmod-update/src/repo"
	"mcmod-update/src/repo/curseforge/v1/schema"
)

var _ repo.FileFetcher = &Adaptor{}

type Adaptor struct {
	apiKey string
}

func NewAdaptor(apiKey string) *Adaptor {
	return &Adaptor{apiKey: apiKey}
}

func (a *Adaptor) GetLatestModFile(modId int32, verison, modLoader string) (*model.File, error) {
	cli := NewClient(a.apiKey, verison, schema.ModLoader(modLoader))
	files, err := cli.GetModFiles(modId, 0, 1)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	file := files[0]

	if !file.IsAvailable {
		return nil, nil
	}

	sha1 := ""
	md5 := ""
	for _, fh := range file.Hashes {
		if fh.Algo == schema.Sha1 {
			sha1 = fh.Value
		}
		if fh.Algo == schema.Md5 {
			md5 = fh.Value
		}
	}
	hash := sha1
	if hash == "" {
		hash = md5
	}

	reqDeps := make([]int32, 0, len(file.Dependencies))
	optDeps := make([]int32, 0, len(file.Dependencies))
	for _, dep := range file.Dependencies {
		if dep.RelationType == schema.RequiredDependency {
			reqDeps = append(reqDeps, dep.ModID)
		}
		if dep.RelationType == schema.OptionalDependency {
			optDeps = append(optDeps, dep.ModID)
		}
	}

	mf := &model.File{
		ModID:        file.ModID,
		FileID:       file.ID,
		ModName:      file.DisplayName,
		DispName:     file.DisplayName,
		FileName:     file.FileName,
		ReleaseType:  file.ReleaseType.String(),
		Hash:         hash,
		Date:         file.FileDate,
		DownloadUrl:  file.DownloadURL,
		GameVersions: file.GameVersions,
		McVersion:    verison,
		ModLoader:    modLoader,
		RequiredDeps: reqDeps,
		OptionalDeps: optDeps,
	}

	return mf, nil
}

func (a *Adaptor) GetLatestModFileWithDeps(modId int32, version,
	modLoader string, optionalDep bool) ([]*model.File, error) {
	file, err := a.GetLatestModFile(modId, version, modLoader)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	files := make([]*model.File, 0, 1+len(file.RequiredDeps)+len(file.OptionalDeps))
	files = append(files, file)

	for _, req := range file.RequiredDeps {
		reqFile, err := a.GetLatestModFile(req, version, modLoader)
		if err != nil {
			return nil, err
		}
		files = append(files, reqFile)
	}

	if optionalDep {
		for _, opt := range file.OptionalDeps {
			optFile, err := a.GetLatestModFile(opt, version, modLoader)
			if err != nil {
				return nil, err
			}
			files = append(files, optFile)
		}
	}

	return files, nil
}
