package crowdin

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Publicly available high level functions generally combining several API calls

const polldelaysec = 5 // Defines delay between each api call when polling a progress status

// Lookup buildId for current project
func (crowdin *Crowdin) GetBuildId() (buildId int, err error) {

	crowdin.log("GetBuildId()")

	var opt ListProjectBuildsOptions
	rl, err := crowdin.ListProjectBuilds(&opt)
	if err != nil {
		return 0, err
	}
	for _, v := range rl.Data {
		if (v.Data.ProjectId == crowdin.config.projectId) && (v.Data.Status == "finished") {
			buildId = v.Data.Id
		}
	}
	if buildId == 0 {
		return 0, errors.New("Can't find a build for this project or build is in progress.")
	}
	return buildId, nil
}

// Lookup projectId
func (crowdin *Crowdin) GetProjectId(projectName string) (projectId int, err error) {

	crowdin.log("GetProjectId()")

	var opt ListProjectsOptions
	rl, err := crowdin.ListProjects(&opt)
	if err != nil {
		return 0, err
	}

	for _, v := range rl.Data {
		if v.Data.Name == projectName {
			projectId = v.Data.ID
		}
	}
	if projectId == 0 {
		return 0, errors.New("Can't find project.")
	}
	return projectId, nil
}

// BuildAllLg - Build a project for all languages
// Options to export:
//   - translated strings only Y/N
//   - approved strings only Y/N
// Update buildProgress
func (crowdin *Crowdin) BuildAllLg(buildTOinSec int, translatedOnly bool, approvedOnly bool) (buildId int, err error) {
	crowdin.log("BuildAllLg()")

	// Invoke build
	var bo BuildProjectTranslationOptions
	// keep bo.BranchId nil
	bo.Languages = nil

	bo.SkipUntranslatedStrings = translatedOnly
	if crowdin.config.apiBaseURL == API_CROWDINDOTCOM {
		bo.ExportApprovedOnly = approvedOnly // crowdin.com
	} else {
		if approvedOnly {
			bo.ExportWithMinApprovalsCount = 1 // Enterprise
		}
	}

	rb, err := crowdin.BuildProjectTranslation(&bo)
	if err != nil {
		return buildId, errors.New("\nBuild Err.")
	}
	buildId = rb.Data.Id
	crowdin.log(fmt.Sprintf("	BuildId=%d", buildId))

	// Poll build status with a timeout
	crowdin.log("	Poll build status crowdin.CheckProjectBuildStatus()")
	timer := time.NewTimer(time.Duration(buildTOinSec) * time.Second)
	defer timer.Stop()
	rp := &ResponseCheckProjectBuildStatus{}
	for rp.Data.Status = rb.Data.Status; rp.Data.Status != "finished" && rp.Data.Status != "canceled"; { // initial value is read from previous API call
		time.Sleep(polldelaysec * time.Second) // delay between each call
		rp, err = crowdin.CheckProjectBuildStatus(&CheckProjectBuildStatusOptions{BuildId: buildId})
		if err != nil {
			crowdin.log(fmt.Sprintf(" Error CheckProjectBuildStatus()=%s", err))
			return buildId, err
			// break
		}
		select {
		case <-timer.C:
			err = errors.New("Build Timeout.")
			return buildId, err
		default:
		}
	}

	if rp.Data.Status != "finished" {
		err = errors.New(fmt.Sprintf("	Build Error:%s", rp.Data.Status))
	}
	return buildId, err
}

// Download a build of the current project
//    outputFileNamePath  required
//    projectId           required if projectName is not provided
//    buildId             optional
// limitation: total number of project directories needs to be 500 max
func (crowdin *Crowdin) DownloadBuild(outputFileNamePath string, buildId int) (err error) {

	// Get URL for downloading
	rd, err := crowdin.DownloadProjectTranslations(&DownloadProjectTranslationsOptions{buildId})
	if err != nil {
		return errors.New("DownloadBuild() - Error getting URL for download.")
	}
	url := rd.Data.Url

	// Actual downloading
	err = crowdin.DownloadFile(url, outputFileNamePath)

	return err
}

// Lookup fileId in current project
//    CrowdinFileName required - full Crowdin path to file (To be noted: does not include the project name)
//		Returns Id and crowdin file name
func (crowdin *Crowdin) LookupFileId(CrowdinFileName string) (id int, name string, err error) {

	crowdin.log(fmt.Sprintf("LookupFileId()\n"))

	// Lookup fileId in Crowdin
	dirId := 0
	crowdinFile := strings.Split(CrowdinFileName, "/")

	crowdin.log(fmt.Sprintf("  len=%d\n", len(crowdinFile)))
	crowdin.log(fmt.Sprintf("  crowdinFile %v\n", crowdinFile))
	// crowdin.log(fmt.Sprintf("  crowdinFile[1] %s\n", crowdinFile[1] ))

	switch l := len(crowdinFile); l {
	case 0:
		return 0, "", errors.New("LookupFileId() - Crowdin file name should not be null.")
	case 1: // no directory so dirId is 0 - value is like "a_file_name"
	case 2: // no directory so dirId is 0 - value is like "/a_file_name"
	default: // l > 1
		// Lookup end directoryId
		// Get a list of all the project folders
		listDirs, err := crowdin.ListDirectories(&ListDirectoriesOptions{Limit: 500})
		if err != nil {
			return 0, "", errors.New("LookupFileId() - Error listing project directories.")
		}
		if len(listDirs.Data) > 0 {
			// Lookup last directory's Id
			dirId = 0
			for i, dirName := range crowdinFile { // Go down the directory branch
				crowdin.log(fmt.Sprintf("  idx %d dirName %s len %d dirId %d", i, dirName, len(crowdinFile), dirId))
				if i > 0 && i < len(crowdinFile)-1 { // 1st entry is empty and we're done once we reach the file name (last item of the slice).
					for _, crwdPrjctDirName := range listDirs.Data { // Look up in list of project dirs the right one
						crowdin.log(fmt.Sprintf("  check -> crwdPrjctDirName.Data.DirectoryId %d crwdPrjctDirName.Data.Name %s", crwdPrjctDirName.Data.DirectoryId, crwdPrjctDirName.Data.Name))
						if crwdPrjctDirName.Data.DirectoryId == dirId && crwdPrjctDirName.Data.Name == dirName {
							dirId = crwdPrjctDirName.Data.Id // Bingo get that Id
							crowdin.log(fmt.Sprintf("  BINGO dirId=%d Crowdin dir name %s", dirId, crwdPrjctDirName.Data.Name))
							break // Done for that one
						}
					}
					if dirId == 0 {
						return 0, "", errors.New(fmt.Sprintf("UpdateFile() - Error: can't match directory names with Crowdin path."))
					}
				}
			}
			if dirId == 0 {
				return 0, "", errors.New(fmt.Sprintf("UpdateFile() - Error: can't match directory names with Crowdin path."))
			}
		} else {
			return 0, "", errors.New("UpdateFile() - Error: mismatch between # of folder found and # of folder expected.")
		}
	}

	crowdinFilename := crowdinFile[len(crowdinFile)-1] // Get file name
	crowdin.log(fmt.Sprintf("  crowdinFilename %s\n", crowdinFilename))

	// Look up file
	listFiles, err := crowdin.ListFiles(&ListFilesOptions{DirectoryId: dirId, Limit: 500})
	if err != nil {
		return 0, "", errors.New("UpdateFile() - Error listing files.")
	}

	fileId := 0
	for _, list := range listFiles.Data {
		crowdin.log(fmt.Sprintf("  check -> list.Data.Name %s", list.Data.Name))
		if list.Data.Name == crowdinFilename {
			fileId = list.Data.Id
			crowdin.log(fmt.Sprintf("  BINGO fileId=%d File name %s", fileId, crowdinFilename))
			break // found it
		}
	}

	if fileId == 0 {
		return 0, "", errors.New(fmt.Sprintf("UpdateFile() - Can't find file %s in Crowdin.", crowdinFilename))
	}

	crowdin.log(fmt.Sprintf("  fileId=%d\n", fileId))
	return fileId, crowdinFilename, nil
}

// Update a file of the current project
//    LocalFileName  required
//    CrowdinFileName required
//    updateOption required needs to be either: clear_translations_and_approvals, keep_translations or keep_translations_and_approvals
//		Returns file Id and rev
func (crowdin *Crowdin) Update(CrowdinFileName string, LocalFileName string, updateOption string) (fileId int, revId int, err error) {

	crowdin.log(fmt.Sprintf("Update()\n"))

	// Lookup fileId in Crowdin
	fileId, crowdinFilename, err := crowdin.LookupFileId(CrowdinFileName)
	if err != nil {
		crowdin.log(fmt.Sprintf("  err=%s\n", err))
		return 0, 0, err
	}

	crowdin.log(fmt.Sprintf("Update() fileId=%d fileName=%s\n", fileId, crowdinFilename))

	// Send local file to storageId
	addStor, err := crowdin.AddStorage(&AddStorageOptions{FileName: LocalFileName})
	if err != nil {
		return 0, 0, errors.New("UpdateFile() - Error adding file to storage.")
	}
	storageId := addStor.Data.Id

	// fmt.Printf("Directory Id = %d, filename= %s, fileId %d storageId= %d\n", dirId, crowdinFilename, fileId, storageId)

	// Update file
	updres, err := crowdin.UpdateFile(fileId, &UpdateFileOptions{StorageId: storageId, UpdateOption: updateOption})

	// Delete storage
	err1 := crowdin.DeleteStorage(&DeleteStorageOptions{StorageId: storageId})

	if err != nil {
		crowdin.log(fmt.Sprintf("UpdateFile() - error updating file %v", updres))
		return 0, 0, errors.New("UpdateFile() - Error updating file.") //
	}

	if err1 != nil {
		crowdin.log(fmt.Sprintf("UpdateFile() - error deleting storage %v", err1))
	}

	revId = updres.Data.RevisionId

	crowdin.log(fmt.Sprintf("UpdateFile() - result %v", updres))

	return fileId, revId, nil
}

// Obtain a list of string Ids for a given file of the current project.
// Use a filter on "identifier" "text" or "context"
// Parameters:
//  - provide path/filename
//	- a filter string (empty mean "all")
//	- filter on "identifier" "text" or "context"
// Returns:
//	- string IDs in a slice of ints if results found
//	- err (nil if no error)
//
func (crowdin *Crowdin) GetStringIDs(fileName string, filter string, filterType string) (list []int, err error) {

	crowdin.log(fmt.Sprintf("GetStringIDs(%s, %s, %s)\n", fileName, filter, filterType))

	// Lookup fileId in Crowdin
	fileId, _, err := crowdin.LookupFileId(fileName)
	if err != nil {
		crowdin.log(fmt.Sprintf("  err=%s\n", err))
		return list, err
	}

	// Get the string IDs
	limit := 500
	opt := ListStringsOptions{
		FileId: fileId,
		Scope:  filterType,
		Filter: filter,
		Limit:  limit,
	}

	// Pull ListStrings as long as it returns data
	for offset := 0; offset < MAX_RESULTS; offset += limit {
		opt.Offset = offset

		res, err := crowdin.ListStrings(&opt)
		if err != nil {
			crowdin.log(fmt.Sprintf("  err=%s\n", err))
			return list, err
		}

		if len(res.Data) <= 0 {
			break
		}

		crowdin.log(fmt.Sprintf(" - Page of results #%d\n", (offset/limit)+1))

		for _, v := range res.Data {
			list = append(list, v.Data.ID) // Add data to slice
		}
	}

	return list, nil
}

type T_UploadTranslationFileParams struct {
	LocalFileName       string // File containing the translations to upload
	CrowdinFileName     string // File in Crowdin where the translations will end up
	LanguageId          string // Langugage ID as per Crowdin spec and defined as target in the project
	ImportEqSuggestions bool   // Defines whether to add translation if it's the same as the source string
	AutoApproveImported bool   // Mark uploaded translations as approved
	TranslateHidden     bool   // Allow translations upload to hidden source strings
	ResponseTimeOut     int    // in seconds. The upload operation can take several minutes.
	// The original TO will be restored after operation finishes (ok or not)
}

// Upload a translation file
// Params:
// 	- File containing the translations to upload
// 	- File in Crowdin where the translations will end up
// 	- Language ID as per Crowdin spec and defined as target in the project
// 	- Defines whether to add translation if it's the same as the source string
// 	- Mark uploaded translations as approved
// 	- Allow translations upload to hidden source strings
// 	- in seconds. The upload operation can take several minutes. 0 means no change.
// 		The original TO will be restored after operation finishes (ok or not)
//	Returns the source fileId (0 if error) and err != nil if error
func (crowdin *Crowdin) UploadTranslationFile(params T_UploadTranslationFileParams) (fileId int, err error) {
	crowdin.log(fmt.Sprintf("UploadTranslationFile(%v)\n", params))

	// Lookup fileId in Crowdin
	fileId, crowdinFilename, err := crowdin.LookupFileId(params.CrowdinFileName)
	if err != nil {
		crowdin.log(fmt.Sprintf("  err=%s\n", err))
		return fileId, err
	}

	crowdin.log(fmt.Sprintf("UploadTranslationFile() fileId=%d fileName=%s\n", fileId, crowdinFilename))

	// Send local file to storageId
	addStor, err := crowdin.AddStorage(&AddStorageOptions{FileName: params.LocalFileName})
	if err != nil {
		crowdin.log(fmt.Sprintf("  Error adding file to storage %s\n", err))
		return fileId, errors.New("UploadTranslationFile() - Error adding file to storage.")
	}
	storageId := addStor.Data.Id

	// fmt.Printf("Directory Id = %d, filename= %s, fileId %d storageId= %d\n", dirId, crowdinFilename, fileId, storageId)

	// Upload file
	if params.ResponseTimeOut > 0 { // If a specific to has been defined
		crowdin.PushTimeouts()                         //  Backup comm timeouts
		crowdin.SetTimeouts(0, params.ResponseTimeOut) // Set new TO for this call
	}
	upldres, err := crowdin.UploadTranslations(params.LanguageId,
		&UploadTranslationsOptions{
			StorageID:           storageId,
			FileID:              fileId,
			ImportEqSuggestions: params.ImportEqSuggestions,
			AutoApproveImported: params.AutoApproveImported,
			TranslateHidden:     params.TranslateHidden,
		})
	if params.ResponseTimeOut > 0 { // If a specific to has been defined
		crowdin.PopTimeouts() // Restore current timeouts
	}

	// Delete storage
	err1 := crowdin.DeleteStorage(&DeleteStorageOptions{StorageId: storageId})

	crowdin.log(fmt.Sprintf("UploadTranslationFile() - uploading %s result %v\n", params.LocalFileName, upldres))
	if err != nil {
		crowdin.log(fmt.Sprintf("UploadTranslationFile() - upload error - %s", err))
		return fileId, errors.New("UploadTranslationFile() - Error uploading file.")
	}

	if err1 != nil {
		// Not a fatal err, just log the error
		crowdin.log(fmt.Sprintf("UploadTranslationFile() - error deleting storage %v", err1))
	}

	return fileId, nil
}


// GetShortLangFileProgress() - Get a simple file completion info for a specific language
//	 Returns a percentage of completion for both translation and approval (0 if error).
func (crowdin *Crowdin) GetShortLangFileProgress(fileId int, langId string) (translationProgress int, approvalProgress int, err error) {
	crowdin.log(fmt.Sprintf("GetShortLangFileProgress()\n"))

	opt := GetFileProgressOptions{FileId: fileId, Limit: 500}
	res, err := crowdin.GetFileProgress(&opt)
	if err == nil {
		// Lookup for language in res
		for _,v := range res.Data {
			if v.Data.LanguageId == langId {
				return v.Data.TranslationProgress, v.Data.ApprovalProgress, nil // found it: done
			}
		}
		crowdin.log(fmt.Sprintf("GetShortLangFileProgress() - language %s not found in %v", langId, res))
		err = errors.New("GetShortLangFileProgress() - Can't find language in result.")
	}
	return 0, 0, err

}
