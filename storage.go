package crowdin

import (
	"encoding/json"
	//"errors"
	"fmt"
	// "io"
	// "net/http"
	// "net/url"
	//"os"
	//"strconv"
	// "time"
	// "github.com/mreiferson/go-httpclient"
)

// ListStorages - List existing storages
// {protocol}://{host}/api/v2/storages
func (crowdin *Crowdin) ListStorages(options *ListStoragesOptions) (*ResponseListStorages, error) {

	crowdin.log("\nListStorages()")

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL + "storages"), body: options})

	if err != nil {
		fmt.Printf("\nREPONSE:%s\n", response)
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseListStorages
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// AddStorage - Add storage API call. Upload a file to a storage space.
// {protocol}://{host}/api/v2/storages
func (crowdin *Crowdin) AddStorage(options *AddStorageOptions) (*ResponseAddStorage, error) {

	crowdin.log("\nAddStorage()")

	// Prepare URL and params
	var p postOptions
	p.urlStr = fmt.Sprintf(crowdin.config.apiBaseURL + "storages")
	p.body = options
	response, err := crowdin.post(&p)

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseAddStorage
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// GetStorage - Read the file name associated to a storageId
// {protocol}://{host}/api/v2/storages/{storageId}
func (crowdin *Crowdin) GetStorage(options *GetStorageOptions) (*ResponseGetStorage, error) {

	crowdin.log("\nGetStorage()")

	response, err := crowdin.get(&getOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"storages/%v", options.StorageId)})

	if err != nil {
		fmt.Printf("\nREPONSE:%s\n", response)
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseGetStorage
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// deleteStorage - Delete a storage
// {protocol}://{host}/api/v2/storages/{storageId}
func (crowdin *Crowdin) DeleteStorage(options *DeleteStorageOptions) error {

	crowdin.log("\nDeleteStorage()")

	response, err := crowdin.del(&delOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"storages/%v", options.StorageId)})

	if err != nil {
		//fmt.Printf("\nREPONSE:%s\n",response)
		crowdin.log(err)
		return err
	}

	crowdin.log(string(response))

	return nil
}
