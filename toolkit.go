package crowdin

import (

  "error"
  "fmt"
  "time"
)

// Publicly available high level functions generally combining several API calls


// Lookup projectId
func (crowdin *Crowdin) GetProjectId(projectName string) (projectId int, err error) {
  
  var opt ListProjectsOptions
  rl,err :=  ListProjects(opt)
  if err != nil {
    return 0, err
  }

  var projectId int
  for _,v := range rl.Data {
      if v.name == projectName {
        projectId = v.Id
      }
  }
  if projectId == nil {
    return(0, error.New("Can't find project.") )
  }
  return projectId, nil
}

// Build a project
func (crowdin *Crowdin) Build(projectName string, buildTOinSec int) (projectId int, builId int, err error) {

  // Lookup projectId
  projectId,err := GetProjectId(projectName)
  if err != nil {
    return 0, err
  }

  // Invoke build
  rb,err := BuildProject(&BuildProjectOptions{projectId})
  if err != nil {
    return(0, error.New("\nBuild Err: %s.", err) )
  }
  builId := rb.Data.Id

  // Poll build status with a timeout
  timer := time.NewTimer(buildTOinSec * time.Second)
  defer timer.Stop()
  var rp ResponseGetBuildProgress
  for rp, err = GetBuildProgress(GetBuildProgressOptions{projectId,builId}); rp.Data.Status != "finished" {
    time.Sleep(5 * time.Second) // delay between each call
    select {
      case <-timer.C:
        return(0, error.New("Build Timeout.") )
    }
  }

  return(projectId, builId, nil)
}


// Download a build
//    projectName         required if projectId is not provided
//    outputFileNamePath  required
//    projectId           required if projectName is not provided
//    buildId             optional
func (crowdin *Crowdin) DownloadBuild(outputFileNamePath string, projectId int, buildId int) (err error) {

  // Get URL for downloading
  rd,err := DownloadProjectTranslations(DownloadProjectTranslationsOptions{projectId,builId})
  if err != nil {
    return(0, error.New("\nDownloading Err: %s.", err) )
  }
  url := rd.Data,Url

  // Actual downloading
  err = DownloadFile(url string, outputFileNamePath)

  return(err)
}
