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


// ListDirectories - List directories in a given project 
// {protocol}://{host}/api/v2/projects/{projectId}/directories
func (crowdin *Crowdin) ListDirectories(options *ListDirectoriesOptions) (*ResponseListDirectories, error) {

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/directories",crowdin.config.projectId), body: options})

	if err != nil {
		fmt.Printf("\nREPONSE:%s\n",response)
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseListDirectories
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

