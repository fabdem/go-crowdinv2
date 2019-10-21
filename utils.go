package crowdin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	//"mime/multipart"
	"net/http"
	"net/http/httputil"

	"os"
	"time"
)

type postOptions struct {
	urlStr string
	body   interface{}
}

type getOptions struct {
	urlStr string
	body   interface{}
}

// params - extra params
// fileNames - key = dir
func (crowdin *Crowdin) post(options *postOptions) ([]byte, error) {

	fmt.Printf("\nCreate Request:")
	fmt.Printf("\nBody: %s", options.body)

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(options.body)
	req, err := http.NewRequest("POST", options.urlStr, buf)

	// var jsonStr = []byte(`{"branchId":null, "targetLanguagesId":[]}`)
	// req, err := http.NewRequest("POST", options.urlStr, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+crowdin.config.token)
	req.Header.Set("Content-Type", "application/json")
	fmt.Printf("\nHeaders: %v\n", req.Header)

	dump, err := httputil.DumpRequestOut(req, true)
	crowdin.log(dump)
	// fmt.Printf("\nDump: %s\n", dump)

	// Run the  request
	response, err := crowdin.config.client.Do(req)
	fmt.Printf("\nDo response: %v\n%v", response, err)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return bodyResponse, APIError{What: fmt.Sprintf("Status code: %v", response.StatusCode)}
	}

	return bodyResponse, nil
}

func (crowdin *Crowdin) get(options *getOptions) ([]byte, error) {

	response, err := crowdin.getResponse(options)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return bodyResponse, APIError{What: fmt.Sprintf("Status code: %v", response.StatusCode)}
	}

	return bodyResponse, nil
}

func (crowdin *Crowdin) getResponse(options *getOptions) (*http.Response, error) {

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(options.body)

	req, err := http.NewRequest("GET", options.urlStr, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+crowdin.config.token)

	fmt.Printf("\nRequest:%v\nError:%v\n", req, err)

	response, err := crowdin.config.client.Do(req)
	fmt.Printf("\nResponse:%v\nError:%v\n", response, err)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func (crowdin *Crowdin) DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (crowdin *Crowdin) log(a interface{}) {
	if crowdin.debug {
		log.Println(a)
		if crowdin.logWriter != nil {
			timestamp := time.Now().Format(time.RFC3339)
			msg := fmt.Sprintf("%v: %v", timestamp, a)
			fmt.Fprintln(crowdin.logWriter, msg)
		}
	}
}

// APIError holds data of errors returned from the API.
type APIError struct {
	What string
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v", e.What)
}
