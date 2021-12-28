package crowdin

import (
	"encoding/json"
	"errors"
	"fmt"
	//"strconv"
)

// ListTranslationApprovals - List translation approvals of a file, a string or a translation
// {protocol}://{host}/api/v2/projects/{projectId}/approvals
//
// 
// 
func (crowdin *Crowdin) ListTranslationApprovals(options *ListTranslationApprovalsOptions) (*ResponseListTranslationApprovals, error) {
	crowdin.log(fmt.Sprintf("ListTranslationApprovals(%d)\n", crowdin.config.projectId))

	translationID	:= options.TranslationID   
	fileID          := options.FileID          
	stringID		:= options.StringID		
	languageID		:= options.LanguageID		

	if !(translationID > 0 || (languageID > 0 && (fileID > 0 || stringID > 0))) { // required
		crowdin.log(fmt.Sprintf("	Error - Minimum nb of parameters not met.\n"))
		return nil, errors.New("insufficient parameters.")
	}

	response, err := crowdin.patch(&patchOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/approvals", crowdin.config.projectId),
		body:   options})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListTranslationApprovals
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil

}
