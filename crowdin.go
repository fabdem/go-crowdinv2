package crowdin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mreiferson/go-httpclient"
)

var (
	// Default value for API URL
	apiBaseURL = "https://crowdin.com/api/v2/"

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
	buildProgress int
	debug         bool
	logWriter     io.Writer
}

// Set connection and read/write timeouts for the subsequent new connections
func SetTimeouts(cnctTOinSecs, rwTOinSecs int) {
	connectionTOinSecs = time.Duration(cnctTOinSecs)
	readwriteTOinSecs = time.Duration(rwTOinSecs)
}

// Read current build progress status from Crowdin structure
// That value is updated when a build is running and GetBuildProgress() polled.
func (crowdin *Crowdin) GetPercentBuildProgress() int {
	return crowdin.buildProgress
}

// New - a create new instance of Crowdin API V2.
func New(token string, projectId int, apiurl string, proxy string) (*Crowdin, error) {

	var proxyUrl *url.URL
	var err error

	if len(apiurl) > 0 { // If a specific URL is defined (Crowdin Enterprise) insert it in the URL
		apiBaseURL = apiurl
	}

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



