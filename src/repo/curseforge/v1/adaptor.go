package v1

import (
	"mcmod-update/src/model"
	"mcmod-update/src/repo/curseforge/v1/schema"
)

type Adaptor struct {
	cli *Client
}

func NewAdaptor(client *Client) *Adaptor {
	return &Adaptor{cli: client}
}

func (a *Adaptor) GetModLatestFile(modId int32) (*model.File, error) {
	files, err := a.cli.GetModFiles(modId, 0, 1)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	file := files[0]

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
		DispName:     file.DisplayName,
		FileName:     file.FileName,
		ReleaseType:  file.ReleaseType.String(),
		Hash:         hash,
		Date:         file.FileDate,
		DownloadUrl:  file.DownloadURL,
		GameVersions: file.GameVersions,
		McVersion:    a.cli.gameVersion,
		ModLoader:    a.cli.modLoader.String(),
		RequiredDeps: reqDeps,
		OptionalDeps: optDeps,
	}

	return mf, nil
}

func (a *Adaptor) GetModLatestFileWithDeps(modId int32, optional bool) ([]*model.File, error) {
	file, err := a.GetModLatestFile(modId)
	if err != nil {
		return nil, err
	}

	files := make([]*model.File, 0, 1+len(file.RequiredDeps)+len(file.OptionalDeps))
	files = append(files, file)

	for _, req := range file.RequiredDeps {
		reqFile, err := a.GetModLatestFile(req)
		if err != nil {
			return nil, err
		}
		files = append(files, reqFile)
	}

	if optional {
		for _, opt := range file.OptionalDeps {
			optFile, err := a.GetModLatestFile(opt)
			if err != nil {
				return nil, err
			}
			files = append(files, optFile)
		}
	}

	return files, nil
}
