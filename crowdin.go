package crowdin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mreiferson/go-httpclient"
	// "go-httpclient"
)

const MAX_RESULTS       = 1000000                       // 1M lines maximum for any api responses
const MAX_RES_PER_PAGE  = 500                           // Max nber of lines per page returned by API calls.
const API_CROWDINDOTCOM = "https://crowdin.com/api/v2/" // url for crowdin.com (non Enterprise version)

const DEFAULT_CONNEXION_TO_S = 5  // seconds
const DEFAULT_RW_TO_S        = 40 // seconds
const DEFAULT_DOWNLOAD_TO_S  = 10 // seconds

var (
	// Default value for API URL
	apiBaseURL = API_CROWDINDOTCOM
	
	// Modifiable (yuck - hence the deprecated fcts at the bottom of file) default values for timeouts in seconds
	DEFAULT_connectionTO = time.Duration(DEFAULT_CONNEXION_TO_S) * time.Second
	DEFAULT_readwriteTO  = time.Duration(DEFAULT_RW_TO_S)        * time.Second
)

// Crowdin API V2 wrapper
type Crowdin struct {
	config struct {
		apiBaseURL          string
		token               string
		projectId           int
		client              *http.Client
		currentConnectionTO time.Duration
		currentReadwriteTO  time.Duration
		currentDownloadTO   time.Duration
		savConnectionTO     time.Duration
		savReadwriteTO      time.Duration
		savDownloadTO      	time.Duration
		proxyUrl            *url.URL
	}
	buildProgress int
	debug         bool
	logWriter     io.Writer
}

// New - a create new instance of Crowdin API V2.
func New(token string, projectId int, apiurl string, proxy string) (*Crowdin, error) {

	var proxyUrl *url.URL
	var err error

	// Default values for timeouts in seconds
	var connectionTO = DEFAULT_connectionTO
	var readwriteTO  = DEFAULT_readwriteTO 

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
		ConnectTimeout:   connectionTO,
		ReadWriteTimeout: readwriteTO,
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
	s.config.currentConnectionTO = connectionTO
	s.config.currentReadwriteTO  = readwriteTO
	s.config.currentDownloadTO	 = time.Duration(DEFAULT_DOWNLOAD_TO_S) * time.Second
	s.config.savConnectionTO 	 = connectionTO
	s.config.savReadwriteTO      = readwriteTO
	s.config.savDownloadTO       = s.config.currentDownloadTO
	s.config.proxyUrl            = proxyUrl

	return s, nil
}

func (crowdin *Crowdin) Close() {
	crowdin.log(fmt.Sprintf("Close() API\n"))
	crowdin = nil
}

// Read current build progress status from Crowdin structure
// That value is updated when a build is running and GetBuildProgress() polled.
func (crowdin *Crowdin) GetPercentBuildProgress() int {
	return crowdin.buildProgress
}

// Set connection and read/write timeouts
//  0 means doesn't change value
func (crowdin *Crowdin) SetSessionTimeouts(cnctTO, rwTO time.Duration) {

	if cnctTO > 0 {
		crowdin.config.currentConnectionTO = cnctTO
	}
	if rwTO > 0 {
		crowdin.config.currentReadwriteTO = rwTO
	}

	transport := &httpclient.Transport{
		ConnectTimeout:   crowdin.config.currentConnectionTO,
		ReadWriteTimeout: crowdin.config.currentReadwriteTO,
		Proxy:            http.ProxyURL(crowdin.config.proxyUrl),
	}
	defer transport.Close()

	crowdin.config.client = &http.Client{
		Transport: transport,
	}
}

// Get connection and read/write timeouts
func (crowdin *Crowdin) GetSessionTimeouts() (cnctTO, rwTO time.Duration) {
	return crowdin.config.currentConnectionTO, crowdin.config.currentReadwriteTO
}

// Set download timeouts
func (crowdin *Crowdin) SetDownloadTimeout(dwnldTO time.Duration) {
	if dwnldTO > time.Duration(0) {
		crowdin.config.currentDownloadTO = dwnldTO
	}
}

// Get download timeout
func (crowdin *Crowdin) GetDownloadTimeouts() (dwnldTO time.Duration) {
	return crowdin.config.currentDownloadTO
}

// Save all timeout values
func (crowdin *Crowdin) PushTimeouts() {
	crowdin.config.savConnectionTO, crowdin.config.savReadwriteTO, crowdin.config.savDownloadTO = crowdin.config.currentConnectionTO, crowdin.config.currentReadwriteTO, crowdin.config.currentDownloadTO
}

// Restore previously saved timeout values
func (crowdin *Crowdin) PopTimeouts() {
	crowdin.SetSessionTimeouts(crowdin.config.savConnectionTO, crowdin.config.savReadwriteTO)
	crowdin.SetDownloadTimeout(crowdin.config.savDownloadTO)	
}

// Reset communication (connection and read/write) timeouts to their default values
func (crowdin *Crowdin) ResetTimeoutsToDefault() {
	crowdin.SetSessionTimeouts(time.Duration(DEFAULT_CONNEXION_TO_S) * time.Second, time.Duration(DEFAULT_RW_TO_S) * time.Second)
	crowdin.SetDownloadTimeout(time.Duration(DEFAULT_DOWNLOAD_TO_S) * time.Second)
}

// SetDebug - traces errors if it's set to true.
func (crowdin *Crowdin) SetDebug(debug bool, logWriter io.Writer) {
	crowdin.debug = debug
	crowdin.logWriter = logWriter
}

// GetDebugWriter - get writer
func (crowdin *Crowdin) GetDebugWriter() (logWriter io.Writer) {
	return(crowdin.logWriter)
}



// DEPRECATED Set connection and read/write timeouts for the subsequent new connections 
func SetDefaultTimeouts(cnctTO, rwTO time.Duration) {
	DEFAULT_connectionTO = cnctTO
	DEFAULT_readwriteTO  = rwTO
}


// DEPRECATED Set connection and read/write timeouts
//  0 means doesn't change value
func (crowdin *Crowdin) SetTimeouts(cnctTO, rwTO time.Duration) {
	crowdin.SetSessionTimeouts(cnctTO, rwTO)
}


// DEPRECATED Get connection and read/write timeouts
func (crowdin *Crowdin) GetTimeouts() (connectionTO, rwTO time.Duration) {
	return crowdin.GetSessionTimeouts()
}
