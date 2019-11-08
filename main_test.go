package main

import (
  "testing"
  "fmt"
  "net/http"
  "github.com/go-redis/redis"
)

const (
  storeId = "50kshop.myshopify.com"
  authToken = "2ba3b6f0d72834e7b643a8e992d06633"
)

func TestGetAllStores(t *testing.T) {
  rdb := redis.NewClient(&redis.Options{
      Addr:     "localhost:6379", // use default Addr
      Password: "",               // no password set
      DB:       2,                // use default DB
  })

  stores := getAllStores(rdb)
  fmt.Println("stores", stores)
}

func TestUpdateScriptTag(t *testing.T) {
  client := &http.Client{}
  scriptTag := getScriptTag(storeId, authToken, client)

  newSrc := "https://cdn.shopify.com/s/files/1/0070/3666/5911/files/whatschat_29db10f4-71f3-4925-a97b-57f3ad1da070.js?1502"

  fmt.Println("updating ", scriptTag.Src, " to ", newSrc)
  updateScriptTag(storeId, authToken, client, scriptTag, newSrc)
}
