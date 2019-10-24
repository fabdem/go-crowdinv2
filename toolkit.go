package crowdin

import (

  "errors"
  "fmt"
  "time"
)

// Publicly available high level functions generally combining several API calls


// Lookup buildId for current project
func (crowdin *Crowdin) GetBuildId() (buildId int, err error) {
  var opt ListProjectBuildsOptions
  rl,err :=  crowdin.ListProjectBuilds(&opt)
  if err != nil {
    return 0, err
  }
  for _,v := range rl.Data {
      if (v.Data.ProjectId == crowdin.config.projectId) && (v.Data.Status == "finished") {
        buildId = v.Data.Id
      }
  }
  if buildId == 0 {
    return 0, errors.New("Can't find a build for this project or build in progress.")
  }
  return buildId, nil
}


// Lookup projectId
func (crowdin *Crowdin) GetProjectId(projectName string) (projectId int, err error) {

  fmt.Printf("")
  var opt ListProjectsOptions
  rl,err :=  crowdin.ListProjects(&opt)
  if err != nil {
    return 0, err
  }

  for _,v := range rl.Data {
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
func (crowdin *Crowdin) BuildAllLg(buildTOinSec int) (buildId int, err error) {
  fmt.Printf("\nIn BuildAllLg()")

  // Invoke build
  var bo BuildProjectOptions
  // bo.ProjectId = crowdin.config.projectId
  bo.BranchId = 0
  bo.Languages = nil
  rb,err := crowdin.BuildProject(&bo)
  if err != nil {
    return buildId, errors.New("\nBuild Err.")
  }
  buildId = rb.Data.Id
  fmt.Printf("\nBuildId=%v", buildId)

  // Poll build status with a timeout
  timer := time.NewTimer(time.Duration(buildTOinSec) * time.Second)
  defer timer.Stop()
  var rp *ResponseGetBuildProgress
  for ;rp.Data.Status != "finished"; {
    time.Sleep(5 * time.Second) // delay between each call
	  rp, err = crowdin.GetBuildProgress(&GetBuildProgressOptions{buildId})
    select {
      case <-timer.C:
        err = errors.New("Build Timeout.")
		break
    }
  }

  return buildId, err
}


// Download a build of the current project
//    outputFileNamePath  required
//    projectId           required if projectName is not provided
//    buildId             optional
func (crowdin *Crowdin) DownloadBuild(outputFileNamePath string, buildId int) (err error) {

  // Get URL for downloading
  rd,err := crowdin.DownloadProjectTranslations(&DownloadProjectTranslationsOptions{buildId})
  if err != nil {
    return errors.New("\nDownloading Err.")
  }
  url := rd.Data.Url

  // Actual downloading
  err = crowdin.DownloadFile(url, outputFileNamePath)

  return err
}
