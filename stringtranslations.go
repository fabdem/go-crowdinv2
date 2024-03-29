package crowdin

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// ListTranslationApprovals - List translation approvals of a file, a string or a translation
// {protocol}://{host}/api/v2/projects/{projectId}/approvals
// 
func (crowdin *Crowdin) ListTranslationApprovals(options *ListTranslationApprovalsOptions) (*ResponseListTranslationApprovals, error) {
	crowdin.log(fmt.Sprintf("ListTranslationApprovals(%d)", crowdin.config.projectId))

	if !(options.TranslationID > 0 || (len(options.LanguageID) > 0 && (options.FileID > 0 || options.StringID > 0))) { // required
		crowdin.log(fmt.Sprintf("	Error - Minimum nb of parameters not met."))
		return nil, errors.New("insufficient parameters.")
	}

	var limit string
	if options.Limit > 0 {
		limit = strconv.Itoa(options.Limit)
	}

	var offset string
	if options.Offset > 0 {
		offset = strconv.Itoa(options.Offset)
	}

	var translationID string
	if options.TranslationID > 0 {
		translationID = strconv.Itoa(options.TranslationID)
	}

	var fileID string
	if options.FileID > 0 {
		fileID = strconv.Itoa(options.FileID)
	}

	var stringID string
	if options.StringID > 0 {
		stringID = strconv.Itoa(options.StringID)
	}

	languageID		:= options.LanguageID		


	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/approvals", crowdin.config.projectId),
		params: map[string]string{
			"fileId":			fileID,
			"stringId":			stringID,
			"languageId":		languageID,
			"translationId": 	translationID,
			"limit":            limit,
			"offset":           offset,
		},
	})

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




// GetTranslation - Get translation details
// {protocol}://{host}/api/v2/projects/{projectId}/translations/{translationId}
// 
func (crowdin *Crowdin) GetTranslation(options *GetTranslationOptions) (*ResponseGetTranslation, error) {
	crowdin.log(fmt.Sprintf("GetTranslation(%d)", options.TranslationID))

	if options.TranslationID <= 0 { // required
		crowdin.log(fmt.Sprintf("	Error - Need a non 0 TransactionID.\n"))
		return nil, errors.New("insufficient parameters.")
	}

	var DenormalizePlaceholders string
	if options.DenormalizePlaceholders > 0 {
		DenormalizePlaceholders = strconv.Itoa(options.DenormalizePlaceholders)
	}

	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/translations/%v", crowdin.config.projectId, options.TranslationID),
		params: map[string]string{
			"denormalizePlaceholders":			DenormalizePlaceholders,
		},
	})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseGetTranslation
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}

// AddApproval - Add Approval API call. Set proofread approval to a translation.
// {protocol}://{host}/api/v2/storages
func (crowdin *Crowdin) AddApproval(options *AddApprovalOptions) (*ResponseAddApproval, error) {
	crowdin.log("AddApproval()")

	// Prepare URL and params
	var p postOptions
	p.urlStr = fmt.Sprintf(crowdin.config.apiBaseURL + "projects/%v/approvals", crowdin.config.projectId) //, options.TranslationID))
	p.body = options
	crowdin.log(fmt.Sprintf("\n	postOptions:%s", p))
	response, err := crowdin.post(&p)
	if err != nil {
		crowdin.log(fmt.Sprintf("\n	post() error:%s\n%s", err, response))
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseAddApproval
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}


// ListLanguageTranslations - List translation from a project for a given language
// Optionnaly crosql can be used to filter the resust set.
// {protocol}://{host}/api/v2/projects/{projectId}/languages/{languageId}/translations
// 
func (crowdin *Crowdin) ListLanguageTranslations(options *ListLanguageTranslationsOptions) (*ResponseListLanguageTranslations, error) {
	crowdin.log(fmt.Sprintf("ListLanguageTranslations(%d - %s)", crowdin.config.projectId, options.LanguageID))

	if len(options.Croql) > 0 && (len(options.StringIDs) > 0 || len(options.LabelIDs) > 0 || options.FileID > 0 ) { // required
		crowdin.log(fmt.Sprintf("	Error - Can't use croql along with other filtering options."))
		return nil, errors.New("Can't use croql along with other filtering options.")
	}

	var limit string
	if options.Limit > 0 {
		limit = strconv.Itoa(options.Limit)
	}

	var offset string
	if options.Offset > 0 {
		offset = strconv.Itoa(options.Offset)
	}

	var fileID string
	if options.FileID > 0 {
		fileID = strconv.Itoa(options.FileID)
	}

	var denormalizePlaceholders string
	if options.FileID > 0 {
		fileID = strconv.Itoa(options.FileID)
	}

	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/languages/%v/translations", crowdin.config.projectId, options.LanguageID),
		params: map[string]string{
			"stringIds":		       options.StringIDs,
			"labelIds":   	           options.LabelIDs,
			"fileId":			       fileID,
			"croql":   	               options.Croql,
			"denormalizePlaceholders": denormalizePlaceholders,
			"limit":                   limit,
			"offset":                  offset,
		},
	})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListLanguageTranslations
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}
