package crowdin

import (
	"encoding/json"
	// "errors"
	"fmt"
	// "strconv"
)

// ListWorkflowsSteps - List workflow steps
// {protocol}://{host}/api/v2/projects/{projectId}/workflow-steps
//
// 
// 
func (crowdin *Crowdin) ListWorkflowsSteps(options *ListWorkflowsStepsOptions) (*ResponseListWorkflowsSteps, error) {
	crowdin.log(fmt.Sprintf("ListWorkflowsSteps(%d)\n", crowdin.config.projectId))


	response, err := crowdin.patch(&patchOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/workflow-steps", crowdin.config.projectId),
		body:   options})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListWorkflowsSteps
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil

}
