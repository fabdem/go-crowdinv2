package crowdin

import (
	"time"
)

// ListProjectBuilds - List Project Build API call
type ListProjectBuildsOptions struct {
	// ProjectId 			int
	// Body struct {
	BranchId int `json:hasManagerAccess,omitempty`
	Limit    int `json:limit,omitempty`
	Offset   int `json:offset,omitempty`
	//}
}


type ResponseListProjectBuilds struct {
	Data []struct {
		Data struct {
			Id         int    `json:"id"`
			ProjectId  int    `json:"projectId"`
			Status     string `json:"status"`
			Progress   int    `json:"progress"`
			Attributes struct {
				BranchId             int   `json:"branchId,omitempty"`
				TargetLanguageIds    []int `json:"targetLanguageIds,omitempty"`
				ExportTranslatedOnly bool  `json:"exportTranslatedOnly"`
				ExportApprovedOnly   bool  `json:"exportApprovedOnly"`
			} `json:"attributes"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// ListProjects - List Projects API call
type ListProjectsOptions struct {
	GroupId          int `json:groupId,omitempty`
	HasManagerAccess int `json:hasManagerAccess,omitempty`
	Limit            int `json:limit,omitempty`
	Offset           int `json:offset,omitempty`
}

type ResponseListProjects struct {
	Data []struct {
		Data struct {
			Id                   int       `json:"id"`
			GroupId              int       `json:"groupId"`
			UserId               int       `json:"userId"`
			SourceLanguageId     int       `json:"sourceLanguageId"`
			TargetLanguageIds    []int     `json:"targetLanguageIds"`
			JoinPolicy           string    `json:"joinPolicy"`
			LanguageAccessPolicy string    `json:"languageAccessPolicy"`
			Type                 int       `json:"type"`
			Name                 string    `json:"name"`
			Cname                string    `json:"cname"`
			Identifier           string    `json:"identifier"`
			Description          string    `json:"description"`
			Visibility           string    `json:"visibility"`
			Logo                 []byte    `json:"logo"`
			Background           string    `json:"background"`
			IsExternal           bool      `json:"isExternal"`
			ExternalType         string    `json:"externalType"`
			AdvancedWorkflowId   int       `json:"advancedWorkflowId"`
			HasCrowdsourcing     bool      `json:"hasCrowdsourcing"`
			CreatedAt            time.Time `json:"createdAt"`
			UpdatedAt            time.Time `json:"updatedAt"`
		} `json:"data"`
	} `json:"data"`
	Pagination []struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// UpdateFile - Update a file API call
type UpdateFileOptions struct {
	StorageId        				int 			`json:"storageId"`
	UpdateOption     				string		`json:"updateOption,omitempty"` //
	ImportOptions struct {
		ContentSegmentation 	bool 			`json:"contentSegmentation,omitempty"`
		TranslateContent 			bool 			`json:"translateContent,omitempty"`
		TranslateAttributes 	bool 			`json:"translateAttributes,omitempty"`
		TranslatableElements  []string  `json:"translatableElements,omitempty"`
	}	 `json:"importOptions,omitempty"`
	ExportOptions struct {
		ContentSegmentation 	bool 			`json:"contentSegmentation,omitempty"`
	}	 `json:"exportOptions,omitempty"`
}

type ResponseUpdateFile struct {
	Data struct {
		Id            int    `json:"id"`
		ProjectId     int    `json:"projectId"`
		BranchId      int    `json:"branchId"`
		DirectoryId   int    `json:"directoryId"`
		Name          string `json:"name"`
		Title         string `json:"title"`
		Type          string `json:"type"`
		RevisionId    int    `json:"revisionId"`
		Status        string `json:"status"`
		Priority      string `json:"priority"`
		ImportOptions struct {
			FirstLineContainsHeader bool `json:"firstLineContainsHeader"`
			ImportTranslations      bool `json:"importTranslations"`
			Scheme                  struct {
				Identifier   int `json:"identifier"`
				SourcePhrase int `json:"sourcePhrase"`
				En           int `json:"en"`
				De           int `json:"de"`
			} `json:"scheme"`
		} `json:"importOptions"`
		ExportOptions struct {
			ExportPattern string `json:"exportPattern"`
		} `json:"exportOptions"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"data"`
}

type ResponseUpdate struct {
	Data struct {
		Id            int    `json:"id"`
		ProjectId     int    `json:"projectId"`
		BranchId      int    `json:"branchId"`
		DirectoryId   int    `json:"directoryId"`
		Name          string `json:"name"`
		Title         string `json:"title"`
		Type          string `json:"type"`
		RevisionId    int    `json:"revisionId"`
		Status        string `json:"status"`
		Priority      string `json:"priority"`
		ImportOptions struct {
			FirstLineContainsHeader bool `json:"firstLineContainsHeader"`
			ImportTranslations      bool `json:"importTranslations"`
			Scheme                  struct {
				Identifier   int `json:"identifier"`
				SourcePhrase int `json:"sourcePhrase"`
				En           int `json:"en"`
				De           int `json:"de"`
			} `json:"scheme"`
		} `json:"importOptions"`
		ExportOptions struct {
			ExportPattern string `json:"exportPattern"`
		} `json:"exportOptions"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"data"`
}

// GetProjectBuildProgressOptions are options for Check Project Build Status api call
type GetBuildProgressOptions struct {
	// Project Identifier.
	// ProjectId int
	BuildId int
}

// {"data":{"id":47,"projectId":17,"status":"inProgress","progress":11,"attributes":{"branchId":null,"targetLanguageIds":[],"exportTranslatedOnly":false,"exportApprovedOnly":false}}}
type ResponseGetBuildProgress struct {
	Data struct {
		Id         int    `json:"id"`
		ProjectId  int    `json:"projectId"`
		Status     string `json:"status"`
		Progress   int    `json:"progress"`
		Attributes struct {
			BranchId             int   `json:"branchId,omitempty"`
			TargetLanguageIds    []int `json:"targetLanguageIds,omitempty"`
			ExportTranslatedOnly bool  `json:"exportTranslatedOnly"`
			ExportApprovedOnly   bool  `json:"exportApprovedOnly"`
		} `json:"attributes"`
	} `json:"data"`
}

// DownloadProjectTranslationsOptions are options for  DownloadProjectTranslations api call
type DownloadProjectTranslationsOptions struct {
	// Project Identifier.
	// ProjectId int
	// Build Identifier.
	BuildId int
}

type ResponseDownloadProjectTranslations struct {
	Data struct {
		Url      string `json:"url"`
		ExpireIn string `json:"expireIn"`
	} `json: "data"`
}

// GetProjectBuilds api call
type ResponseGetProjectBuilds struct {
	Data []struct {
		Data struct {
			Id         int    `json:"id"`
			ProjectId  int    `json:"projectId"`
			BranchId   int    `json:"branchId"`
			LanguageId []int  `json:"languageId"`
			Status     string `json:"status"`
			Progress   struct {
				Percent           int `json:"percent"`
				CurrentLanguageId int `json:"currentLanguageId"`
				CurrentFileId     int `json:"currentFileId"`
			} `json:"progress"`
		} `json:"data"`
	} `json:"data"`
	Pagination []struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// GetLanguageProgress api call
type ResponseGetLanguageProgress struct {
	Data []struct {
		Data struct {
			LanguageId                int `json:"languageId"`
			PhrasesCount              int `json:"phrasesCount"`
			PhrasesTranslatedCount    int `json:"phrasesTranslatedCount"`
			PhrasesApprovedCount      int `json:"phrasesApprovedCount"`
			PhrasesTranslatedProgress int `json:"phrasesTranslatedProgress"`
			PhrasesApprovedProgress   int `json:"phrasesApprovedProgress"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// BuildProjectOptions are options for BuildProject api call
type BuildProjectOptions struct {
	// ProjectId int		 // Project Identifier.
	// Body      struct {
	BranchId int `json:"branchId,omitempty"` // Branch Identifier. - optional
	// Specify target languages for build.
	// Leave this field empty to build all target languages
	Languages []string `json:"targetLanguagesId,omitempty"`
	// }
}

type ResponseBuildProject struct {
	Data struct {
		Id          int      `json:"id"`
		ProjectId   int      `json:"projectId"`
		BranchId    int      `json:"branchId"`
		LanguagesId []string `json:"languagesId"`
		Status      string   `json:"status"`
		Progress    struct {
			Percent           int `json:"percent"`
			CurrentLanguageId int `json:"currentLanguageId"`
			CurrentFileId     int `json:"currentFileId"`
		} `json:"progress"`
	} `json:"data"`
}

// ListStoragesOptions are options for ListStorages api call
type ListStoragesOptions struct {
	Limit  int `json:"limit,omitempty"`  // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset int `json:"offset,omitempty"` // Offset in collection - optional
}

// ResponseListStorages are response for ListStorages api call
type ResponseListStorages struct {
	Data []struct {
		Data struct {
			Id       int    `json:"id"`
			FileName string `json:"fileName"`
		} `json:"data"`
	} `json:"data"`
	Pagination []struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// AddStorageOptions are options for AddStorage api call
type AddStorageOptions struct {
	FileName string `json:"filename"` // Filename and path of hte file to upload to storage
}

// ResponseAddStorage are response for AddStorage api call
type ResponseAddStorage struct {
	Data struct {
		Id       int    `json:"id"`
		FileName string `json:"fileName"`
	} `json:"data"`
}

// GetStorageOptions are options for GetStorage api call
type GetStorageOptions struct {
	StorageId int `json:"storageid"`
}

// ResponseGetStorage are response for ListStorages api call
type ResponseGetStorage struct {
	Data struct {
		Id       int    `json:"id"`
		FileName string `json:"fileName"`
	} `json:"data"`
}

// DelStorageOptions are options for DelStorage api call
type DeleteStorageOptions struct {
	StorageId int `json:"storageid"`
}

// ListDirectoriesOptions are options for ListDirectories api call
type ListDirectoriesOptions struct {
	BranchId    int `json:"branchId,omitempty"`
	DirectoryId int `json:"directoryId,omitempty"`
	Recursion   int `json:"recursion,omitempty"`
	Limit       int `json:"limit,omitempty"`  // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset      int `json:"offset,omitempty"` // Offset in collection - optional
}

// ResponseListDirectories are response for ListDirectories api call
type ResponseListDirectories struct {
	Data []struct {
		Data struct {
			Id            int       `json:"id"`
			ProjectId     int       `json:"projectId"`
			BranchId      int       `json:"branchId"`
			DirectoryId   int       `json:"directoryId"`  // Actually parentId
			Name          string    `json:"name"`
			Title         string    `json:"title"`
			ExportPattern string    `json:"exportPattern"`
			Priority      string    `json:"priority"`
			CreatedAt     time.Time `json:"createdAt"`
			UpdatedAt     time.Time `json:"updatedAt"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// ListFilesOptions are options for ListFiles api call
type ListFilesOptions struct {
	BranchId    int `json:"branchId,omitempty"`
	DirectoryId int `json:"directoryId,omitempty"`
	Recursion   int `json:"recursion,omitempty"`
	Limit       int `json:"limit,omitempty"`  // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset      int `json:"offset,omitempty"` // Offset in collection - optional
}

// ResponseListFiles are response for ListFiles api call
type ResponseListFiles struct {
	Data []struct {
		Data struct {
			Id          int    `json:"id"`
			ProjectId   int    `json:"projectId"`
			BranchId    int    `json:"branchId"`
			DirectoryId int    `json:"directoryId"`
			Name        string `json:"name"`
			Title       string `json:"title"`
			Type        string `json:"type"`
			RevisionId  int    `json:"revisionId"`
			Status      string `json:"status"`
			Priority    string `json:"priority"`
			Attributes  struct {
				MimeType string `json:"mimeType"`
				FileSize int    `json:"fileSize"`
			} `json:"attributes"`
			ImportOptions struct {
				FirstLineContainsHeader bool `json:"firstLineContainsHeader"`
				ImportTranslations      bool `json:"importTranslations"`
				Scheme                  struct {
					Identifier   int `json:"identifier"`
					SourcePhrase int `json:"sourcePhrase"`
					En           int `json:"en"`
					De           int `json:"de"`
				} `json:"scheme"`
			} `json:"importOptions"`
			ExportOptions struct {
				ExportPattern string `json:"exportPattern"`
			} `json:"exportOptions"`
			ExportPattern string    `json:"exportPattern"`
			CreatedAt     time.Time `json:"createdAt"`
			UpdatedAt     time.Time `json:"updatedAt"`
			Revision      int       `json:"revision"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

type responseGeneral struct {
	Success bool `json:"success"`
}


// ListFilesOptions are options for ListFileRevisions api call
type ListFileRevisionsOptions struct {
	Limit       int `json:"limit,omitempty"`  // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset      int `json:"offset,omitempty"` // Offset in collection - optional
}

// ResponseListFiles are response for ListFiles api call
type ResponseListFileRevisions struct {
	Data []struct {
		Data struct {
			Id                int `json:"id"`
			ProjectId         int `json:"projectId"`
			FileId            int `json:"fileId"`
			RestoreToRevision int `json:"restoreToRevision"`
			Info              struct {
				Added struct {
					Strings int `json:"strings"`
					Words   int `json:"words"`
				} `json:"added"`
				Deleted struct {
					Strings int `json:"strings"`
					Words   int `json:"words"`
				} `json:"deleted"`
				Updated struct {
					Strings int `json:"strings"`
					Words   int `json:"words"`
				} `json:"updated"`
			} `json:"info"`
			Date time.Time `json:"date"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}
