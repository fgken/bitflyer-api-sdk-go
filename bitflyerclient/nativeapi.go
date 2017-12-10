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

type Page struct {
	Count  uint64
	Before uint64
	After  uint64
}

type ProductCode string

const (
	BTC_JPY    ProductCode = "BTC_JPY"
	FX_BTC_JPY ProductCode = "FX_BTC_JPY"
)

func addPage(values url.Values, page Page) url.Values {
	values.Add("Count", strconv.FormatUint(page.Count, 10))
	values.Add("Before", strconv.FormatUint(page.Before, 10))
	values.Add("After", strconv.FormatUint(page.After, 10))
	return values
}

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
 *  Trade API
 * ==============================
 */

type ResponseGetExecutions struct {
	Id                        uint64
	Child_order_id            string
	Side                      string
	Price                     float64
	Size                      float64
	Commission                float64
	Exec_date                 BitflyerTime
	Child_order_acceptance_id string
}

func (client *Client) GetExecutions(code ProductCode, page Page) ([]ResponseGetExecutions, error) {
	var param requestParam
	param.path = "/v1/me/getexecutions"
	param.method = http.MethodGet
	param.isPrivate = true
	queries := url.Values{}
	queries.Add("product_code", string(code))
	queries = addPage(queries, page)
	param.queryString = queries.Encode()

	body, err := client.do(param)
	if err != nil {
		return nil, err
	}

	result := make([]ResponseGetExecutions, 0)
	if err := json.Unmarshal(*body, &result); err != nil {
		log.Printf("error: %v\n", err)
	}

	return result, err
}
