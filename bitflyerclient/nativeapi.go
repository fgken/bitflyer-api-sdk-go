package bitflyerclient

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	//"github.com/k0kubun/pp"
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

func (page *Pagenation) init() {
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
 *  Public API
 * ==============================
 */

/* --- Get Order Book (Board) --- */
type BoardOrder struct {
	Price float64
	Size  float64
}

type GetBoardResponse struct {
	Mid_price float64
	Bids      []BoardOrder
	Asks      []BoardOrder
}

func (client *Client) GetBoard() (*GetBoardResponse, error) {
	reqParam := requestParam{
		path:      "/v1/getboard",
		method:    http.MethodGet,
		isPrivate: false,
	}
	queries := url.Values{}
	queries.Add("product_code", string(client.productCode))
	reqParam.queryString = queries.Encode()

	respBody, err := client.do(reqParam)
	if err != nil {
		return nil, err
	}

	var result GetBoardResponse
	if err := json.Unmarshal(*respBody, &result); err != nil {
		log.Printf("error: %v\n", err)
	}

	return &result, err
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
	Id                        int64
	Child_order_id            string
	Side                      string
	Price                     float64
	Size                      float64
	Commission                float64
	Exec_date                 BitflyerTime
	Child_order_acceptance_id string
}

func (client *Client) GetExecutions(param *GetExecutionsParam) ([]GetExecutionsResponse, error) {
	reqParam := requestParam{
		path:      "/v1/me/getexecutions",
		method:    http.MethodGet,
		isPrivate: true,
	}
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

/* --- Get Child Orders --- */
type GetChildOrdersParam struct {
	Page Pagenation
    Product_code string
    Child_order_state string
    Child_order_id string
    Child_order_acceptance_id string
    Parent_order_id string
}

func NewGetChildOrdersParam() *GetChildOrdersParam {
	var param GetChildOrdersParam
	param.Page.init()
	return &param
}

type GetChildOrdersResponse struct {
	Id                        int64
	Child_order_id            string
    Product_code              string
    Child_order_type          string
	Side                      string
	Price                     float64
    Average_price             float64
	Size                      float64
    Child_order_state         string
    Expire_date               BitflyerTime
    Child_order_date          BitflyerTime
	Child_order_acceptance_id string
    Outstanding_size          float64
    Cancel_size               float64
    Executed_size             float64
    Total_commission          float64
}

func (client *Client) GetChildOrders(param *GetChildOrdersParam) ([]GetChildOrdersResponse, error) {
	reqParam := requestParam{
		path:      "/v1/me/getchildorders",
		method:    http.MethodGet,
		isPrivate: true,
	}
	queries := url.Values{}
	queries.Add("product_code", string(client.productCode))
	queries = addPagenation(queries, param.Page)
    if param.Child_order_state != "" {
        queries.Add("child_order_state", param.Child_order_state)
    }
    if param.Child_order_id != "" {
        queries.Add("child_order_id", param.Child_order_id)
    }
    if param.Child_order_acceptance_id != "" {
        queries.Add("child_order_acceptance_id", param.Child_order_acceptance_id)
    }
    if param.Parent_order_id != "" {
        queries.Add("parent_order_id", param.Parent_order_id)
    }
	reqParam.queryString = queries.Encode()

	respBody, err := client.do(reqParam)
	if err != nil {
		return nil, err
	}

	result := make([]GetChildOrdersResponse, 0)
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

/* Submit New Parent Order (Special Order) */
const (
	SIMPLE = "SIMPLE"
	IFD    = "IFD"
	OCO    = "OCO"
	IFDOCO = "IFDOCO"
)

type ParentOrder struct {
	Product_code   string  `json:"product_code"`
	Condition_type string  `json:"condition_type"`
	Side           string  `json:"side"`
	Size           float64 `json:"size"`
	Price          float64 `json:"price"`
	Trigger_price  float64 `json:"trigger_price"`
	Offset         float64 `json:"offset"`
}

type SendParentOrderParam struct {
	Order_method     string        `json:"order_method"`
	Minute_to_expire uint64        `json:"minute_to_expire"`
	Time_in_force    string        `json:"time_in_force"`
	Parameters       []ParentOrder `json:"parameters"`
}

func NewSendParentOrderParam() *SendParentOrderParam {
	var param SendParentOrderParam
	param.Minute_to_expire = 43200
	param.Time_in_force = "GTC"
	param.Parameters = make([]ParentOrder, 0)
	return &param
}

type SendParentOrderResponse struct {
	Parent_order_acceptance_id string
}

func (client *Client) SendParentOrder(param *SendParentOrderParam) (*SendParentOrderResponse, error) {
	for i, _ := range param.Parameters {
		param.Parameters[i].Product_code = client.productCode
	}
	reqParam := requestParam{
		path:        "/v1/me/sendparentorder",
		method:      http.MethodPost,
		isPrivate:   true,
		queryString: "",
	}

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

	var result SendParentOrderResponse
	if err := json.Unmarshal(*respBody, &result); err != nil {
		log.Printf("error: %v\n", err)
	}

	return &result, err
}

/* --- List Parent Orders --- */
type GetParentOrdersParam struct {
	Product_code       string
	Page               Pagenation
	Parent_order_state string
}

func NewGetParentOrdersParam() *GetParentOrdersParam {
	var param GetParentOrdersParam
	param.Page.init()
	return &param
}

type GetParentOrdersResponse struct {
	Id                         int64
	Parent_order_id            string
	Product_code               string
	Side                       string
	Parent_order_type          string
	Price                      float64
	Size                       float64
	Parent_order_state         string
	Expire_date                BitflyerTime
	Parent_order_date          BitflyerTime
	Parent_order_acceptance_id string
	Outstanding_size           float64
	Cancel_size                float64
	Executed_size              float64
	total_commission           float64
}

func (client *Client) GetParentOrders(param *GetParentOrdersParam) ([]GetParentOrdersResponse, error) {
	reqParam := requestParam{
		path:      "/v1/me/getparentorders",
		method:    http.MethodGet,
		isPrivate: true,
	}
	queries := url.Values{}
	queries.Add("product_code", string(client.productCode))
	queries = addPagenation(queries, param.Page)
	if param.Parent_order_state != "" {
		queries.Add("parent_order_state", param.Parent_order_state)
	}
	reqParam.queryString = queries.Encode()

	respBody, err := client.do(reqParam)
	if err != nil {
		return nil, err
	}

	result := make([]GetParentOrdersResponse, 0)
	if err := json.Unmarshal(*respBody, &result); err != nil {
		log.Printf("error: %v\n", err)
	}

	return result, err
}

/* --- Get Parent Order Detail --- */
type GetParentOrderParam struct {
	Parent_order_acceptance_id string
	Parent_order_id            string /* should not use */
}

func NewGetParentOrderParam() *GetParentOrderParam {
	var param GetParentOrderParam
	return &param
}

type GetParentOrderResponse struct {
	Id                         int64
	Parent_order_acceptance_id string
	Parent_order_id            string /* should not use */
	Order_method               string
	Minute_to_expire           uint64
	Parameters                 []ParentOrder
}

func (client *Client) GetParentOrder(param *GetParentOrderParam) (*GetParentOrderResponse, error) {
	reqParam := requestParam{
		path:      "/v1/me/getparentorder",
		method:    http.MethodGet,
		isPrivate: true,
	}
	queries := url.Values{}
	queries.Add("parent_order_acceptance_id", param.Parent_order_acceptance_id)
	reqParam.queryString = queries.Encode()

	respBody, err := client.do(reqParam)
	if err != nil {
		return nil, err
	}

	var result GetParentOrderResponse
	if err := json.Unmarshal(*respBody, &result); err != nil {
		log.Printf("error: %v\n", err)
	}

	return &result, err
}
