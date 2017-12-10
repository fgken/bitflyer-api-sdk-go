package bitflyerclient

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	APIEndpointBase = "https://api.bitflyer.jp"
)

const (
	BTC_JPY    = "BTC_JPY"
	FX_BTC_JPY = "FX_BTC_JPY"
)

type Client struct {
	apiKey       string
	apiSecret    string
	endpointBase string
	httpClient   *http.Client
	productCode  string
}

func New(apiKey, apiSecret string) (*Client, error) {
	c := &Client{
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		endpointBase: APIEndpointBase,
		httpClient:   http.DefaultClient,
		productCode:  FX_BTC_JPY,
	}
	return c, nil
}

type requestParam struct {
	path        string
	method      string
	isPrivate   bool
	queryString string
	body        string
}

func (client *Client) do(param requestParam) (*[]byte, error) {
	path := param.path
	if param.queryString != "" {
		path += "?" + param.queryString
	}
	url := client.endpointBase + path

	req, err := http.NewRequest(param.method, url, strings.NewReader(param.body))
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if param.isPrivate {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		text := timestamp + param.method + path + param.body
		mac := hmac.New(sha256.New, []byte(client.apiSecret))
		mac.Write([]byte(text))
		sign := hex.EncodeToString(mac.Sum(nil))

		req.Header.Set("ACCESS-KEY", client.apiKey)
		req.Header.Set("ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("ACCESS-SIGN", sign)
	}

	log.Printf("debug: Send request: %v %v %v\n", url, param.method, param.body)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("error: %v(%v) %v", resp.Status, resp.StatusCode, string(respBody))
		return nil, err
	}

	log.Printf("debug: Receive Response: %v\n", string(respBody))

	return &respBody, nil
}
