package crowdin

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	// "reflect"

)

// ListMTs - List the Machine Translation engines defined in the organization.
//   Optionnaly filter by group. !!!! However, as of 2023-08-01 Crowdin doesn't return the inherited MTs !!!
//                 (Which is annoying if you're trying to determine what MTs are available to a given project.)
// 
// {protocol}://{host}/api/v2/mts
//
func (crowdin *Crowdin) ListMTs(options *ListMTsOptions) (*ResponseListMTs, error) {
	crowdin.log(fmt.Sprintf("ListMTs() Pulls list of machine translation engines from groupId=%d", options.GroupID))

	var groupId string
	if options.GroupID > 0 {
		groupId = strconv.Itoa(options.GroupID)
	}

	var limit string
	if options.Limit > 0 {
		limit = strconv.Itoa(options.Limit)
	}

	var offset string
	if options.Offset > 0 {
		offset = strconv.Itoa(options.Offset)
	}

	var params map[string]string
	if crowdin.config.apiBaseURL == API_CROWDINDOTCOM { // crowdin.com version
		params = map[string]string{
				"limit":   limit,
				"offset":  offset,
			}
	} else {											// Enterprise version
		params = map[string]string{
				"groupId": groupId,
				"limit":   limit,
				"offset":  offset,
			}
	}
	
	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL + "mts"),
		params: params,
	})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListMTs
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}



// TranslateViaMT - Add Translate via MT call.
// {protocol}://{host}/api/v2/mts/{mtId}/translations
func (crowdin *Crowdin) TranslateViaMT(mtID int, options *TranslateViaMTOptions) (*ResponseTranslateViaMT, []string, error) {
	crowdin.log(fmt.Sprintf("TranslateViaMT(%d)", mtID))

	// Prepare URL and params
	var p postOptions
	p.urlStr = fmt.Sprintf(crowdin.config.apiBaseURL + "mts/%v/translations", mtID)
	p.body = options
	crowdin.log(fmt.Sprintf("\n	postOptions:%s", p))
	response, err := crowdin.post(&p)
	if err != nil {
		crowdin.log(fmt.Sprintf("\n	post() error:%s\n%s", err, response))
		return nil, nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseTranslateViaMT
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, nil, err
	}
	
	// Parse the translations - they come in 2 different flavors.
	// Either list of objects (crowdin) or array of strings (Amazon, Google)
	var translations []string
	jsonObj := responseAPI.Data.Translations
	switch obj := jsonObj.(type) {
	case []interface{}:
		for _, v := range obj {
			val, ok := v.(string)
			if !ok {
				return nil, nil, errors.New(fmt.Sprintf("\n	Response - error unmarshalling translations array %s", obj))
			}
			translations = append(translations, val)
		}
	case map[string]interface{}:
		translations = make([]string, len(obj))
		for i := 0; i < len(obj); i++ {
			translations[i] = fmt.Sprintf("%s", obj[strconv.Itoa(i)])
		}
	default:	// Anything else: error
		return nil, nil, errors.New(fmt.Sprintf("\n	Response - error unmarshalling transaltions. Unrecognized type %s - %s", obj, jsonObj))
	}	

	return &responseAPI, translations, nil
}
