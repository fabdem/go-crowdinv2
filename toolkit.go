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
			projectId = v.Data.Id
		}
	}
	if projectId == 0 {
		return 0, errors.New("Can't find project.")
	}
	return projectId, nil
}

// BuildAllLg - Build a project for all languages
// Update buildProgress
func (crowdin *Crowdin) BuildAllLg(buildTOinSec int) (buildId int, err error) {
	crowdin.log("BuildAllLg()")

	// Invoke build
	var bo BuildProjectTranslationOptions
	// bo.ProjectId = crowdin.config.projectId
	bo.BranchId = 0
	bo.Languages = nil
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
			break
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
//    crowdinFileNamePath required - full Crowdin path to file
//		Returns Id and crowdin file name
func (crowdin *Crowdin) LookupFileId(crowdinFileNamePath string) (id int, name string, err error) {

	crowdin.log(fmt.Sprintf("LookupFileId()\n"))

	// Lookup fileId in Crowdin
	dirId := 0
	crowdinFile := strings.Split(crowdinFileNamePath, "/")

	switch l := len(crowdinFile); l {
	case 0:
		return 0, "", errors.New("LookupFileId() - Crowdin file name should not be null.")
	case 1: // no directory so dirId is 0
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
				if i > 1 && i < len(crowdinFile) - 1 { // Skip project name and we're done once we reach the file name (last item of the slice).
					for _, crwdPrjctDirName := range listDirs.Data { // Look up in list of project dirs the right one
						crowdin.log(fmt.Sprintf("  check -> crwdPrjctDirName.Data.DirectoryId %d crwdPrjctDirName.Data.Name %s", crwdPrjctDirName.Data.DirectoryId, crwdPrjctDirName.Data.Name))
						if crwdPrjctDirName.Data.DirectoryId == dirId && crwdPrjctDirName.Data.Name == dirName {
							dirId = crwdPrjctDirName.Data.Id  // Bingo get that Id
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

	crowdinFilename := crowdinFile[len(crowdinFile) - 1]   // Get file name
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
			break   // found it
		}
	}

	if fileId == 0 {
		return 0, "", errors.New(fmt.Sprintf("UpdateFile() - Can't find file %s in Crowdin.", crowdinFilename))
	}
	
	crowdin.log(fmt.Sprintf("  fileId=%d\n", fileId))
	return fileId, crowdinFilename, nil
}

// Update a file of the current project
//    localFileNamePath  required
//    crowdinFileNamePath required
//    updateOption required needs to be either: clear_translations_and_approvals, keep_translations or keep_translations_and_approvals
//		Returns file Id
func (crowdin *Crowdin) Update(crowdinFileNamePath string, localFileNamePath string, updateOption string) (id int, err error) {

	crowdin.log(fmt.Sprintf("Update()\n"))

	// Lookup fileId in Crowdin
	fileId, crowdinFilename, err := crowdin.LookupFileId(crowdinFileNamePath)
	if err != nil {
		return 0, errors.New("UpdateFile() - Can't find file in Crowdin.")
	}

	crowdin.log(fmt.Sprintf("Update() fileId=%d fileName=%s\n", fileId, crowdinFilename))

	// Send local file to storageId
	addStor, err := crowdin.AddStorage(&AddStorageOptions{FileName: localFileNamePath})
	if err != nil {
		return 0, errors.New("UpdateFile() - Error adding file to storage.")
	}
	storageId := addStor.Data.Id

	// fmt.Printf("Directory Id = %d, filename= %s, fileId %d storageId= %d\n", dirId, crowdinFilename, fileId, storageId)

	// Update file
	updres, err := crowdin.UpdateFile(fileId, &UpdateFileOptions{StorageId: storageId, UpdateOption: updateOption})

	// Delete storage
	err1 := crowdin.DeleteStorage(&DeleteStorageOptions{StorageId: storageId})

	if err != nil {
		crowdin.log(fmt.Sprintf("UpdateFile() - error updating file %v", updres))
		return 0, errors.New("UpdateFile() - Error updating file.") //
	}

	if err1 != nil {
		crowdin.log(fmt.Sprintf("UpdateFile() - error deleting storage %v", err1))
	}

	crowdin.log(fmt.Sprintf("UpdateFile() - result %v", updres))

	return fileId, nil
}
