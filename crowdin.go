package crowdin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mreiferson/go-httpclient"
)

const MAX_RESULTS = 1000000  // 1M lines
const API_CROWDINDOTCOM = "https://crowdin.com/api/v2/"  // url for crowdin.com (non Enterprise version)

const DEFAULT_CONNEXION_TO = 5	// seconds
const DEFAULT_RW_TO				= 40 	// seconds

var (
	// Default value for API URL
	apiBaseURL = API_CROWDINDOTCOM

	// Default values for timeouts in seconds
	connectionTOinSecs time.Duration = DEFAULT_CONNEXION_TO
	readwriteTOinSecs  time.Duration = DEFAULT_RW_TO
)


// Crowdin API V2 wrapper
type Crowdin struct {
	config struct {
		apiBaseURL					string
		token								string
		projectId						int
		client							*http.Client
		currentConnectionTO	int
		currentReadwriteTO 	int
		savConnectionTO			int
		savReadwriteTO 			int
		proxyUrl						*url.URL
	}
	buildProgress int
	debug         bool
	logWriter     io.Writer
}

// Set connection and read/write timeouts for the subsequent new connections
// func SetTimeouts(cnctTOinSecs, rwTOinSecs int) {
// 	connectionTOinSecs = time.Duration(cnctTOinSecs)
// 	readwriteTOinSecs = time.Duration(rwTOinSecs)
// }


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
		ConnectTimeout:   DEFAULT_CONNEXION_TO * time.Second,
		ReadWriteTimeout: DEFAULT_RW_TO * time.Second,
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
	s.config.currentConnectionTO	= DEFAULT_CONNEXION_TO
	s.config.currentReadwriteTO		= DEFAULT_RW_TO
	s.config.savConnectionTO			= DEFAULT_CONNEXION_TO
	s.config.savReadwriteTO				= DEFAULT_RW_TO
	s.config.proxyUrl 						= proxyUrl

	return s, nil
}


// Set connection and read/write timeouts
//  0 means doesn't change value
func (crowdin *Crowdin) SetTimeouts(connectionTO, rwTO int) {

	if connectionTO > 0 { crowdin.config.currentConnectionTO	= connectionTO }
	if rwTO > 0 { crowdin.config.currentReadwriteTO 	= rwTO }

	transport := &httpclient.Transport{
		ConnectTimeout:   time.Duration(crowdin.config.currentConnectionTO) * time.Second,
		ReadWriteTimeout: time.Duration(crowdin.config.currentReadwriteTO)  * time.Second,
		Proxy:            http.ProxyURL(crowdin.config.proxyUrl),
	}
	defer transport.Close()

	crowdin.config.client = &http.Client{
		Transport: transport,
		}
}

// Get connection and read/write timeouts
func (crowdin *Crowdin) GetTimeouts()(connectionTO, rwTO int) {
	return crowdin.config.currentConnectionTO, crowdin.config.currentReadwriteTO
}

// Save current timeout values
func (crowdin *Crowdin) PushTimeouts() {
	crowdin.config.savConnectionTO	= crowdin.config.currentConnectionTO
	crowdin.config.savReadwriteTO 	= crowdin.config.currentReadwriteTO
}

// Restore saved timeout values
func (crowdin *Crowdin) PopTimeouts() {
	crowdin.SetTimeouts(crowdin.config.savConnectionTO, crowdin.config.savReadwriteTO)
}

// Reset communication timeouts to their default values
func (crowdin *Crowdin) ResetTimeoutsToDefault() {
	crowdin.SetTimeouts(DEFAULT_CONNEXION_TO, DEFAULT_RW_TO)
}


// SetDebug - traces errors if it's set to true.
func (crowdin *Crowdin) SetDebug(debug bool, logWriter io.Writer) {
	crowdin.debug = debug
	crowdin.logWriter = logWriter
}
