package bitflyerclient

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/k0kubun/pp"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	APIEndpointBase = "https://api.bitflyer.jp"
)

type methodType string

const (
	GET  methodType = "GET"
	POST            = "POST"
)

type httpAPI struct {
	method    methodType
	path      string
	isPrivate bool
}

var apiGetExecution = httpAPI{method: GET, path: "/v1/me/getexecutions", isPrivate: true}

type ResponseGetExecutions struct {
	Id             uint64
	Child_order_id string
	Side           string
	Price          float64
	Size           float64
	Commission     float64
	Exec_date      string
	//Exec_date time.Time
	Child_order_acceptance_id string
}

type Client struct {
	apiKey       string
	apiSecret    string
	endpointBase string
	httpClient   *http.Client
}

type Page struct {
	Count  uint64
	Before uint64
	After  uint64
}

func New(apiKey, apiSecret string) (*Client, error) {
	c := &Client{
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		endpointBase: APIEndpointBase,
		httpClient:   http.DefaultClient,
	}
	return c, nil
}

func (client *Client) do(api httpAPI, query url.Values) (*http.Response, error) {
	method := string(api.method)
	apipath := api.path
	if 0 < len(query) {
		apipath += "?" + query.Encode()
	}
	url := client.endpointBase + apipath
	body := ""

	pp.Println(method, url)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if api.isPrivate {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		text := timestamp + method + apipath + body
		mac := hmac.New(sha256.New, []byte(client.apiSecret))
		mac.Write([]byte(text))
		sign := hex.EncodeToString(mac.Sum(nil))

		pp.Println(client.apiKey)
		req.Header.Set("ACCESS-KEY", client.apiKey)
		req.Header.Set("ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("ACCESS-SIGN", sign)
	}

	return client.httpClient.Do(req)
}

func (client *Client) GetExecutions(page Page) (*[]ResponseGetExecutions, error) {
	queries := url.Values{}
	queries.Add("product_code", "FX_BTC_JPY")
	if page.Count != 0 {
		queries.Add("Count", strconv.FormatUint(page.Count, 10))
	}
	if page.Before != 0 {
		queries.Add("Before", strconv.FormatUint(page.Before, 10))
	}
	if page.After != 0 {
		queries.Add("After", strconv.FormatUint(page.After, 10))
	}

	resp, err := client.do(apiGetExecution, queries)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Println(string(body))
	result := make([]ResponseGetExecutions, 0)
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		log.Println(err)
	}

	pp.Println(result)
	return &result, err
}
