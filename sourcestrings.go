package crowdin

import (
	"encoding/json"
	// "errors"
	"fmt"
	// "io"
	// "net/http"
	// "net/url"
	//"os"
	"strconv"
	// "time"
	// "github.com/mreiferson/go-httpclient"
)

// ListStrings - List Source Strings
// {protocol}://{host}/api/v2/projects/{projectId}/strings
func (crowdin *Crowdin) ListStrings(options *ListStringsOptions) (*ResponseListStrings, error) {

	crowdin.log(fmt.Sprintf("ListDirectories()\n"))

	var fileId string
	if options.FileId > 0 {
		fileId = strconv.Itoa(options.FileId)
	}

	var denormalizePlaceholders string
	if options.DenormalizePlaceholders > 0 {
		denormalizePlaceholders = strconv.Itoa(options.DenormalizePlaceholders)
	}

	var limit string
	if options.Limit > 0 {
		limit = strconv.Itoa(options.Limit)
	}

	var offset string
	if options.Offset > 0 {
		offset = strconv.Itoa(options.Offset)
	}

	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/strings", crowdin.config.projectId),
		params: map[string]string{
			"fileId":                  fileId,
			"denormalizePlaceholders": denormalizePlaceholders,
			"labelIds":                options.LabelIds,
			"filter":                  options.Filter,
			"scope":                   options.Scope,
			"limit":                   limit,
			"offset":                  offset,
		},
	})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListStrings
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}
