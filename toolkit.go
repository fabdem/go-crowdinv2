package crowdin

import (

  "errors"
  "fmt"
  "time"
)

// Publicly available high level functions generally combining several API calls


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

// Build a project
func (crowdin *Crowdin) Build(projectName string, buildTOinSec int) (projectId int, buildId int, err error) {

  // Lookup projectId
  projectId,err = crowdin.GetProjectId(projectName)
  if err != nil {
    return 0,0, err
  }

  // Invoke build
  var bo BuildProjectOptions
  bo.ProjectId = projectId
  bo.Body.BranchId = 0
  bo.Body.Languages = nil
  rb,err := crowdin.BuildProject(&bo)
  if err != nil {
    return projectId, buildId, errors.New("\nBuild Err.")
  }
  buildId = rb.Data.Id

  // Poll build status with a timeout
  timer := time.NewTimer(time.Duration(buildTOinSec) * time.Second)
  defer timer.Stop()
  var rp *ResponseGetBuildProgress
  for ;rp.Data.Status != "finished"; {
    time.Sleep(5 * time.Second) // delay between each call
	rp, err = crowdin.GetBuildProgress(&GetBuildProgressOptions{projectId,buildId})
    select {
      case <-timer.C:
        err = errors.New("Build Timeout.") 
		break
    }
  }

  return projectId, buildId, err
}


// Download a build
//    projectName         required if projectId is not provided
//    outputFileNamePath  required
//    projectId           required if projectName is not provided
//    buildId             optional
func (crowdin *Crowdin) DownloadBuild(outputFileNamePath string, projectId int, buildId int) (err error) {

  // Get URL for downloading
  rd,err := crowdin.DownloadProjectTranslations(&DownloadProjectTranslationsOptions{projectId,buildId})
  if err != nil {
    return errors.New("\nDownloading Err.") 
  }
  url := rd.Data.Url

  // Actual downloading
  err = crowdin.DownloadFile(url, outputFileNamePath)

  return err
}
