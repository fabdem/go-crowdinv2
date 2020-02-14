package crowdin

import (
	"encoding/json"
	//"errors"
	"fmt"
	// "io"
	// "net/http"
	// "net/url"
	//"os"
	//"strconv"
	// "time"
	// "github.com/mreiferson/go-httpclient"
)

// ListFiles - List directories in a given project
// {protocol}://{host}/api/v2/projects/{projectId}/files
func (crowdin *Crowdin) ListDirectories(options *ListDirectoriesOptions) (*ResponseListDirectories, error) {

	crowdin.log(fmt.Sprintf("\nListDirectories()"))

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/directories", crowdin.config.projectId), body: options})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListDirectories
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}

// ListFiles - List files in a given project
// {protocol}://{host}/api/v2/projects/{projectId}/files
func (crowdin *Crowdin) ListFiles(options *ListFilesOptions) (*ResponseListFiles, error) {

	crowdin.log(fmt.Sprintf("\nListFiles()"))

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/files", crowdin.config.projectId), body: options})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListFiles
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}


// UpdateFile - Update a specific file
// {protocol}://{host}/api/v2/projects/{projectId}/files/{fileId}
func (crowdin *Crowdin) UpdateFile(fileId int, options *UpdateFileOptions) (*ResponseUpdateFile, error) {

	crowdin.log(fmt.Sprintf("\nUpdateFile()"))

	response, err := crowdin.put(&putOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/files/%v", crowdin.config.projectId, fileId), body: options})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseUpdateFile
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}
