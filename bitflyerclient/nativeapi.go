package bitflyerclient

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/* ==============================
 *  Common struct and functions
 * ==============================
 */

/* --- Pagenation --- */
type Pagenation struct {
	Count  int64
	Before int64
	After  int64
}

func (page Pagenation) init() {
	page.Count = -1
	page.Before = -1
	page.After = -1
}

func addPagenation(values url.Values, page Pagenation) url.Values {
	if 0 <= page.Count {
		values.Add("count", strconv.FormatInt(page.Count, 10))
	}
	if 0 <= page.Before {
		values.Add("before", strconv.FormatInt(page.Before, 10))
	}
	if 0 <= page.After {
		values.Add("after", strconv.FormatInt(page.After, 10))
	}
	return values
}

/* --- Parse Bitflyer's time format --- */
type BitflyerTime struct {
	time.Time
}

const bitflyerTimeLayout = "2006-01-02T15:04:05.999"

func (bt *BitflyerTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	bt.Time, err = time.Parse(bitflyerTimeLayout, s)
	return err
}

/* ==============================
 *  Trading API
 * ==============================
 */

/* --- Get Execution History --- */
type GetExecutionsParam struct {
	Page Pagenation
}

func NewGetExecutionsParam() *GetExecutionsParam {
	var param GetExecutionsParam
	param.Page.init()
	return &param
}

type GetExecutionsResponse struct {
	Id                        uint64
	Child_order_id            string
	Side                      string
	Price                     float64
	Size                      float64
	Commission                float64
	Exec_date                 BitflyerTime
	Child_order_acceptance_id string
}

func (client *Client) GetExecutions(param *GetExecutionsParam) ([]GetExecutionsResponse, error) {
	var reqParam requestParam
	reqParam.path = "/v1/me/getexecutions"
	reqParam.method = http.MethodGet
	reqParam.isPrivate = true
	queries := url.Values{}
	queries.Add("product_code", string(client.productCode))
	queries = addPagenation(queries, param.Page)
	reqParam.queryString = queries.Encode()

	respBody, err := client.do(reqParam)
	if err != nil {
		return nil, err
	}

	result := make([]GetExecutionsResponse, 0)
	if err := json.Unmarshal(*respBody, &result); err != nil {
		log.Printf("error: %v\n", err)
	}

	return result, err
}

/* --- Send a New Order --- */
type SendChildOrderParam struct {
	Product_code     string  `json:"product_code"`
	Child_order_type string  `json:"child_order_type"`
	Side             string  `json:"side"`
	Price            float64 `json:"price"`
	Size             float64 `json:"size"`
	Minute_to_expire uint64  `json:"minute_to_expire"`
	Time_in_force    string  `json:"time_in_force"`
}

func NewSendChildOrderParam() *SendChildOrderParam {
	var param SendChildOrderParam
	param.Minute_to_expire = 43200
	param.Time_in_force = "GTC"
	return &param
}

type SendChildOrderResponse struct {
	Child_order_acceptance_id string
}

func (client *Client) SendChildOrder(param *SendChildOrderParam) (*SendChildOrderResponse, error) {
	param.Product_code = client.productCode
	var reqParam requestParam
	reqParam.path = "/v1/me/sendchildorder"
	reqParam.method = http.MethodPost
	reqParam.isPrivate = true
	reqParam.queryString = ""

	bodyJson, err := json.Marshal(param)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}

	reqParam.body = string(bodyJson)
	respBody, err := client.do(reqParam)
	if err != nil {
		return nil, err
	}

	var result SendChildOrderResponse
	if err := json.Unmarshal(*respBody, &result); err != nil {
		log.Printf("error: %v\n", err)
	}

	return &result, err
}
