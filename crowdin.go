package crowdin

import (
	"encoding/json"
	//"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	//"os"
	//"strconv"
	"time"

	"github.com/mreiferson/go-httpclient"
)

var (
	apiBaseURL = "https://valve.crowdin.com/api/v2/"

	// Default values for timeouts in seconds
	connectionTOinSecs time.Duration = 5
	readwriteTOinSecs  time.Duration = 40
)

// Crowdin API V2 wrapper
type Crowdin struct {
	config struct {
		apiBaseURL string
		token      string
		projectId  int
		client     *http.Client
	}
	buildProgress  	int
	debug     		bool
	logWriter 		io.Writer
}

// Set connection and read/write timeouts for the subsequent new connections
func SetTimeouts(cnctTOinSecs, rwTOinSecs int) {
	connectionTOinSecs = time.Duration(cnctTOinSecs)
	readwriteTOinSecs = time.Duration(rwTOinSecs)
}

// Read current build status
func (crowdin *Crowdin) ReadBuildProgress() int {
	return crowdin.buildProgress
}


// New - a create new instance of Crowdin API V2.
func New(token string, projectId int, proxy string) (*Crowdin, error) {

	var proxyUrl *url.URL
	var err error

	if len(proxy) > 0 { // If a proxy is defined
		proxyUrl, err = url.Parse(proxy)
		if err != nil {
			fmt.Println("Bad proxy URL", err)
			return nil, err
		}
	}

	transport := &httpclient.Transport{
		ConnectTimeout:   connectionTOinSecs * time.Second,
		ReadWriteTimeout: readwriteTOinSecs * time.Second,
		Proxy:            http.ProxyURL(proxyUrl),
	}
	defer transport.Close()

	s := &Crowdin{}
	s.config.apiBaseURL = apiBaseURL
	s.config.token = token
	s.config.projectId = projectId
	s.config.client = &http.Client{
		Transport: transport,
	}
	return s, nil
}

// SetProject - set project details
//func (crowdin *Crowdin) SetProject(token string, project string) *Crowdin {
//	crowdin.config.token = token
//	crowdin.config.project = project
//	return crowdin
//}

// SetDebug - traces errors if it's set to true.
func (crowdin *Crowdin) SetDebug(debug bool, logWriter io.Writer) {
	crowdin.debug = debug
	crowdin.logWriter = logWriter
}

// ListProjectBuilds - List Project Builds API call. List the project builds
// {protocol}://{host}/api/v2/projects/{projectId}/translations/builds
func (crowdin *Crowdin) ListProjectBuilds(options *ListProjectBuildsOptions) (*ResponseListProjectBuilds, error) {

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/translations/builds", crowdin.config.projectId), body: options})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseListProjectBuilds
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// ListProjects - List projects API call. List the projects and their respective details (incl.Id.)
// {protocol}://{host}/api/v2/projects
func (crowdin *Crowdin) ListProjects(options *ListProjectsOptions) (*ResponseListProjects, error) {

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL + "projects")})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseListProjects
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// DownloadProjectTranslations - Download Project Translations api call
// {protocol}://{host}/api/v2/projects/{projectId}/translations/builds/{buildId}/download
func (crowdin *Crowdin) DownloadProjectTranslations(options *DownloadProjectTranslationsOptions) (*ResponseDownloadProjectTranslations, error) {

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/translations/builds/%v/download", crowdin.config.projectId, options.BuildId)})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseDownloadProjectTranslations
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// GetProjectBuilds - List Project Builds api call
// {protocol}://{host}/api/v2/projects/{projectId}/translations/builds
func (crowdin *Crowdin) GetProjectBuilds() (*ResponseGetProjectBuilds, error) {

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/translations/builds", crowdin.config.projectId)})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseGetProjectBuilds
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// GetBuildProgress - Check Project Build Status api call
// {protocol}://{host}/api/v2/projects/{projectId}/translations/builds/{buildId}
func (crowdin *Crowdin) GetBuildProgress(options *GetBuildProgressOptions) (*ResponseGetBuildProgress, error) {

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/translations/builds/%v", crowdin.config.projectId, options.BuildId)})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseGetBuildProgress
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// GetLanguageProgress() - Get progress info per language API call
// {protocol}://{host}/api/v2/projects/{projectId}/languages/progress
func (crowdin *Crowdin) GetLanguageProgress() (*ResponseGetLanguageProgress, error) {

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/languages/progress", crowdin.config.projectId)})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseGetLanguageProgress
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// BuildProject - Build a project
// {protocol}://{host}/api/v2/projects/{ProjectId}/translations/builds
func (crowdin *Crowdin) BuildProject(options *BuildProjectOptions) (*ResponseBuildProject, error) {

	// Prepare URL and params
	var p postOptions
	p.urlStr = fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/translations/builds", crowdin.config.projectId)
	p.body = options
	response, err := crowdin.post(&p)

	if err != nil {
		fmt.Printf("\ncrowdinV2 - result = %s \n err=", response, err)
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseBuildProject
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
	}

	return &responseAPI, nil
}
