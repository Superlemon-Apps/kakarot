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
}

func newUpdateScriptJob(httpClient *http.Client, storeId string, authToken string, contains string, newSrc string) *UpdateScriptJob{
    return &UpdateScriptJob{
      httpClient: httpClient,
      storeId: storeId,
      authToken: authToken,
      contains: contains,
      newSrc: newSrc,
    }
}

func(self *UpdateScriptJob) Name() string {
  return self.storeId
}

func(self *UpdateScriptJob) Execute() error {
  currentScriptTag := getScriptTag(self.storeId, self.authToken, self.httpClient)
  if currentScriptTag != nil && strings.Contains(currentScriptTag.Src, self.contains){
    fmt.Println("updating ", currentScriptTag.Src, " to ", self.newSrc, " for store ", self.storeId)
    // updateScriptTag(self.storeId, self.authToken, self.httpClient, currentScriptTag, self.newSrc)
  }
  return nil
}
