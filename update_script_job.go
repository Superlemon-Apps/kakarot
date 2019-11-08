package main

import (
  "net/http"
  "fmt"
  "strings"
)

type UpdateScriptJob struct {
  storeId string
  authToken string
  contains string
  newSrc string
  httpClient *http.Client
  Id int
}

func newUpdateScriptJob(httpClient *http.Client, storeId string, authToken string, contains string, newSrc string, id int) *UpdateScriptJob{
    return &UpdateScriptJob{
      httpClient: httpClient,
      storeId: storeId,
      authToken: authToken,
      contains: contains,
      newSrc: newSrc,
      Id: id,
    }
}

func(self *UpdateScriptJob) Name() string {
  return self.storeId
}

func(self *UpdateScriptJob) Execute() error {
  currentScriptTag := getScriptTag(self.storeId, self.authToken, self.httpClient)
  if currentScriptTag != nil && strings.Contains(currentScriptTag.Src, self.contains){
    fmt.Println("updating ", currentScriptTag.Src, " to ", self.newSrc, " for store ", self.storeId, " job no : ", self.Id)
    updateScriptTag(self.storeId, self.authToken, self.httpClient, currentScriptTag, self.newSrc, 3)
  }
  return nil
}
