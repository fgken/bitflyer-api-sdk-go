package main

import (
	"log"
	"os"
	"time"
	"strconv"
	"io/ioutil"
	"net/http"
	"encoding/hex"
	"crypto/hmac"
	"crypto/sha256"
)

var ENDPOINT_URL="https://api.bitflyer.jp/"

func main() {
	api_key := os.Getenv("API_KEY")
	api_secret := os.Getenv("API_SECRET")

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	method := "GET"
	path := "/v1/me/getbalance"
	body := ""

	text := timestamp + method + path + body

	mac := hmac.New(sha256.New, []byte(api_secret))
	mac.Write([]byte(text))
	sign := hex.EncodeToString(mac.Sum(nil))

	req, err := http.NewRequest(method, ENDPOINT_URL + path, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("ACCESS-KEY", api_key)
	req.Header.Set("ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("ACCESS-SIGN", sign)
	//req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	bodya, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bodya))
}
