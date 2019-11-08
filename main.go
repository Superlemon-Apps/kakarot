package main

import (
  "flag"
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/go-redis/redis"
  "bytes"
  "errors"
  "github.com/sankalpjonn/wrq"
)

var (
  contains string
  new_src string
  redis_addr string
)

type ScriptTag struct {
  Id int `json:"id"`
  Src string `json:"src"`
}

type ScriptTagAPIResponse struct {
  ScriptTags []ScriptTag `json:"script_tags"`
}

func getAllStores(rdb *redis.Client) map[string]string {
  var stores = map[string]string{}
  cmd := rdb.Keys("shop:*:settings")
  keys, err := cmd.Result()
  if err != nil {
    panic(err)
  }
  for _, key := range(keys) {
    val := rdb.HMGet(key, "id", "auth_token").Val()
    if len(val) == 2 && val[0] != nil && val[1] != nil{
        stores[val[0].(string)] = val[1].(string)
    }
  }
  return stores
}

func getScriptTag(storeId string, authToken string, httpClient *http.Client) *ScriptTag {
  url := fmt.Sprintf("https://%s/admin/api/2019-10/script_tags.json", storeId)
  req, _ := http.NewRequest("GET", url, nil)
  req.Header.Set("x-shopify-access-token", authToken)
  res, err := httpClient.Do(req)
  if err != nil {
    panic(err)
  }
  defer res.Body.Close()
  if res.StatusCode  == 200 {
    var scriptTags ScriptTagAPIResponse
    json.NewDecoder(res.Body).Decode(&scriptTags)
    return &(scriptTags.ScriptTags[0])
  } else {
    return nil;
  }
}

func updateScriptTag(storeId string, authToken string, httpClient *http.Client, scriptTag *ScriptTag, newSrc string) {
  newScriptTag := ScriptTag{Id: scriptTag.Id, Src:newSrc}
  url := fmt.Sprintf("https://%s/admin/api/2019-10/script_tags/%d.json", storeId, scriptTag.Id)
  requestBody, _ := json.Marshal(map[string]interface{}{"script_tag": newScriptTag})
  req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
  req.Header.Set("x-shopify-access-token", authToken)
  req.Header.Set("content-type", "application/json")
  res, _ := httpClient.Do(req)
  defer res.Body.Close()

  if res.StatusCode != 200 {
    panic(errors.New("found a non 200 status code on updating scripttag"))
  }
}

func main() {
  flag.StringVar(&contains, "contains", "", "updates scripts that contains this string")
	flag.StringVar(&new_src, "new-src", "", "New scripttag url")
  flag.StringVar(&redis_addr, "redis-addr", "localhost:6379", "Redis db address")
	flag.Parse()

  rdb := redis.NewClient(&redis.Options{
      Addr:     redis_addr, // use default Addr
      Password: "",               // no password set
      DB:       2,                // use default DB
  })

  w := wrq.NewWithSettings("shopify", 100, 10)
  defer w.Stop()

  httpClient := &http.Client{}

  stores := getAllStores(rdb)

  for storeId, authToken := range stores {
    job := newUpdateScriptJob(httpClient, storeId, authToken, contains, new_src)
    w.AddJob(job)
  }
}
