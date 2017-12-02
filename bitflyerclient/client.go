package bitflyerclient

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/k0kubun/pp"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

type OrderType string

const (
	MARKET OrderType = "MARKET"
	LIMIT  OrderType = "LIMIT"
)

type OrderSide string

const (
	BUY  OrderSide = "BUY"
	SELL           = "SELL"
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

var apiSendChildOrder = httpAPI{method: POST, path: "/v1/me/sendchildorder", isPrivate: true}

type RequestSendChildOrder struct {
	Product_code     string  `json:"product_code"`
	Child_order_type string  `json:"child_order_type"`
	Side             string  `json:"side"`
	Price            uint64  `json:"price"`
	Size             float64 `json:"size"`
	Minute_to_expire uint64  `json:"minute_to_expire"`
	Time_in_force    string  `json:"time_in_force"`
}

type ResponseSendChildOrder struct {
    Child_order_acceptance_id   string
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

func (client *Client) do(api httpAPI, query url.Values, body string) (*http.Response, error) {
	method := string(api.method)
	apipath := api.path
	if 0 < len(query) {
		apipath += "?" + query.Encode()
	}
	url := client.endpointBase + apipath

	pp.Println(method, url)
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if api.isPrivate {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		text := timestamp + method + apipath + body
		mac := hmac.New(sha256.New, []byte(client.apiSecret))
		mac.Write([]byte(text))
		sign := hex.EncodeToString(mac.Sum(nil))

		req.Header.Set("ACCESS-KEY", client.apiKey)
		req.Header.Set("ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("ACCESS-SIGN", sign)
	}

	//pp.Println(req)

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

	resp, err := client.do(apiGetExecution, queries, "")
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

func (client *Client) SendChildOrder(orderType OrderType, orderSide OrderSide, size float64, price uint64) (*ResponseSendChildOrder, error) {
	bodyParam := RequestSendChildOrder{Minute_to_expire: 43200, Time_in_force: "GTC"}
	bodyParam.Product_code = "FX_BTC_JPY"
	bodyParam.Child_order_type = string(orderType)
	bodyParam.Side = string(orderSide)
	bodyParam.Price = price
	bodyParam.Size = size

	bodyJson, err := json.Marshal(bodyParam)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	queries := url.Values{}
	fmt.Println(string(bodyJson))
	resp, err := client.do(apiSendChildOrder, queries, string(bodyJson))
	if err != nil {
		log.Println(err)
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("error: %s\n", resp.Status))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Println(string(body))
	var result ResponseSendChildOrder
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println(err)
	}

	pp.Println(result)
	return &result, err
}
