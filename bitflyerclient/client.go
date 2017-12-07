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
	BUY       OrderSide = "BUY"
	SELL                = "SELL"
	SIDE_NONE           = ""
)

func oppositeSide(side OrderSide) OrderSide {
	switch side {
	case BUY:
		return SELL
	case SELL:
		return BUY
	}
	panic(fmt.Sprintf("Unexpedted OrderSide: %v", side))
}

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
	Child_order_acceptance_id string
}

var apiSendParentOrder = httpAPI{method: POST, path: "/v1/me/sendparentorder", isPrivate: true}

type ParentOrder struct {
	Product_code   string  `json:"product_code"`
	Condition_type string  `json:"condition_type"`
	Side           string  `json:"side"`
	Size           float64 `json:"size"`
	Price          float64 `json:"price"`
	Trigger_price  float64 `json:"trigger_price"`
	Offset         uint64  `json:"offset"`
}

type RequestSendParentOrder struct {
	Order_method     string        `json:"order_method"`
	Minute_to_expire uint64        `json:"minute_to_expire"`
	Time_in_force    string        `json:"time_in_force"`
	Parameters       []ParentOrder `json:"parameters"`
}

type ResponseSendParentOrder struct {
	Parent_order_acceptance_id string
}

var apiGetBoard = httpAPI{method: GET, path: "/v1/getboard", isPrivate: false}

type BoardOrder struct {
	Price float64
	Size  float64
}

type ResponseGetBoard struct {
	Mid_price float64
	Bids      []BoardOrder
	Asks      []BoardOrder
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

	//pp.Println(method, url)
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

func (client *Client) GetExecutionsById(acceptanceId string) (*[]ResponseGetExecutions, error) {
	queries := url.Values{}
	queries.Add("product_code", "FX_BTC_JPY")
	queries.Add("child_order_acceptance_id", acceptanceId)

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

	//fmt.Println(string(body))
	result := make([]ResponseGetExecutions, 0)
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		log.Println(err)
	}

	pp.Println(result)
	return &result, err
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

	//fmt.Println(string(body))
	result := make([]ResponseGetExecutions, 0)
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		log.Println(err)
	}

	pp.Println(result)
	return &result, err
}

type ResponseGetChildOrders struct {
	Id                        uint64
	Child_order_id            string
	Product_code              string
	Side                      string
	Child_order_type          string
	Price                     float64
	Average_price             float64
	Size                      float64
	Child_order_state         string
	Expire_date               string
	Child_order_date          string
//	Expire_date               time.Time
//	Child_order_date          time.Time
	Child_order_acceptance_id string
	Outstanding_size          float64
	Cancel_size               float64
	Executed_size             float64
	Total_commission          float64
}

var apiGetChildOrders = httpAPI{method: GET, path: "/v1/me/getchildorders", isPrivate: true}

func (client *Client) GetChildOrdersByParentId(parentId string) (*[]ResponseGetChildOrders, error) {
	queries := url.Values{}
	queries.Add("product_code", "FX_BTC_JPY")
	queries.Add("parent_order_id", parentId)

	resp, err := client.do(apiGetChildOrders, queries, "")
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

	//fmt.Println(string(body))
	result := make([]ResponseGetChildOrders, 0)
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		log.Println(err)
	}

	//pp.Println(result)
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
	//fmt.Println(string(bodyJson))
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

	//fmt.Println(string(body))
	var result ResponseSendChildOrder
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println(err)
	}

	pp.Println(result)
	return &result, err
}

func (client *Client) SendParentOrder_IFDOCO(orderSide OrderSide, size float64, entryPrice, limitPrice, stopPrice float64) (*ResponseSendParentOrder, error) {
	var bodyParam RequestSendParentOrder
	bodyParam.Order_method = "IFDOCO"
	bodyParam.Minute_to_expire = 43200
	bodyParam.Time_in_force = "GTC"

	var order ParentOrder

	/* ENTRY(LIMIT) */
	order.Product_code = "FX_BTC_JPY"
	order.Condition_type = "LIMIT"
	order.Side = string(orderSide)
	order.Price = entryPrice
	order.Size = size
	bodyParam.Parameters = append(bodyParam.Parameters, order)

	/* LIMIT ORDER */
	order.Condition_type = "LIMIT"
	order.Side = string(oppositeSide(orderSide))
	order.Price = limitPrice
	bodyParam.Parameters = append(bodyParam.Parameters, order)

	/* STOP ORDER */
	order.Condition_type = "STOP"
	order.Side = string(oppositeSide(orderSide))
	order.Price = 0
	order.Trigger_price = stopPrice
	bodyParam.Parameters = append(bodyParam.Parameters, order)

	bodyJson, err := json.Marshal(bodyParam)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//fmt.Println(string(bodyJson))
	resp, err := client.do(apiSendParentOrder, url.Values{}, string(bodyJson))
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

	var result ResponseSendParentOrder
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println(err)
	}

	return &result, err
}

func (client *Client) GetBoard() (*ResponseGetBoard, error) {
	queries := url.Values{}
	queries.Add("product_code", "FX_BTC_JPY")

	resp, err := client.do(apiGetBoard, queries, "")
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

	//pp.Println(string(body))
	var result ResponseGetBoard
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		log.Println(err)
		return nil, err
	}

	return &result, err
}
