package schema

import "time"

// Reference: https://docs.curseforge.com/#tocS_FileRelationType
type FileReleaseType int

func (r FileReleaseType) String() string {
	switch int(r) {
	case int(Release):
		return "Release"
	case int(Beta):
		return "Beta"
	case int(Alpha):
		return "Alpha"
	default:
		return "Unknown"
	}
}

const (
	Release FileReleaseType = 1
	Beta    FileReleaseType = 2
	Alpha   FileReleaseType = 3
)

// Reference: https://docs.curseforge.com/#tocS_FileStatus
type FileStatusType int

const (
	Processing         FileStatusType = 1
	ChangesRequired    FileStatusType = 2
	UnderReview        FileStatusType = 3
	Approved           FileStatusType = 4
	Rejected           FileStatusType = 5
	MalwareDetected    FileStatusType = 6
	Deleted            FileStatusType = 7
	Archived           FileStatusType = 8
	Testing            FileStatusType = 9
	Released           FileStatusType = 10
	ReadyForReview     FileStatusType = 11
	Deprecated         FileStatusType = 12
	Baking             FileStatusType = 13
	AwaitingPublishing FileStatusType = 14
	FailedPublishing   FileStatusType = 15
)

type FileHash struct {
	Value string   `json:"value"`
	Algo  HashAlgo `json:"algo"`
}

// Reference: https://docs.curseforge.com/#tocS_FileRelationType
type FileRelationType int

const (
	EmbeddedLibrary    FileRelationType = 1
	OptionalDependency FileRelationType = 2
	RequiredDependency FileRelationType = 3
	Tool               FileRelationType = 4
	Incompatible       FileRelationType = 5
	Include            FileRelationType = 6
)

type FileDependency struct {
	ModID        int32            `json:"modId"`
	RelationType FileRelationType `json:"relationType"`
}

// Reference: https://docs.curseforge.com/#tocS_Get%20Mod%20Files%20Response
type GetModFilesResponse struct {
	Data       []*File    `json:"data"`       // The response data
	Pagination Pagination `json:"pagination"` // The response pagination information
}

// Reference: https://docs.curseforge.com/#tocS_File
type File struct {
	ID                   int32                 `json:"id"`                   // The file id
	GameID               int32                 `json:"gameId"`               // The game id related to the mod that this file belongs to
	ModID                int32                 `json:"modId"`                // The mod id
	IsAvailable          bool                  `json:"isAvailable"`          // Whether the file is available to download
	DisplayName          string                `json:"displayName"`          // Display name of the file
	FileName             string                `json:"fileName"`             // Exact file name
	ReleaseType          FileReleaseType       `json:"releaseType"`          // The file release type
	FileStatus           FileStatusType        `json:"fileStatus"`           // Status of the file
	Hashes               []FileHash            `json:"hashes"`               // The file hash (i.e. md5 or sha1)
	FileDate             time.Time             `json:"fileDate"`             // The file timestamp
	FileLength           int64                 `json:"fileLength"`           // The file length in bytes
	DownloadCount        int64                 `json:"downloadCount"`        // The number of downloads for the file
	DownloadURL          string                `json:"downloadUrl"`          // The file download URL
	GameVersions         []string              `json:"gameVersions"`         // List of game versions this file is relevant for
	SortableGameVersions []SortableGameVersion `json:"sortableGameVersions"` // Metadata used for sorting by game versions
	Dependencies         []FileDependency      `json:"dependencies"`         // List of dependencies files
}
