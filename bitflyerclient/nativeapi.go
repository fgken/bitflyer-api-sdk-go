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

/* --- Page Nation --- */
type Page struct {
	count  int64
	before int64
	after  int64
}

func NewPage() *Page {
	return &Page{
		count:  -1,
		before: -1,
		after:  -1,
	}
}

func (page *Page) SetCount(count int64) {
	page.count = count
}

func (page *Page) SetBefore(before int64) {
	page.before = before
}

func (page *Page) SetAfter(after int64) {
	page.after = after
}

func addPage(values url.Values, page *Page) url.Values {
	if 0 <= page.count {
		values.Add("count", strconv.FormatInt(page.count, 10))
	}
	if 0 <= page.before {
		values.Add("before", strconv.FormatInt(page.before, 10))
	}
	if 0 <= page.after {
		values.Add("after", strconv.FormatInt(page.after, 10))
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
 *  Trade API
 * ==============================
 */

/* --- Execution History --- */
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

func (client *Client) GetExecutions(page *Page) ([]ResponseGetExecutions, error) {
	var param requestParam
	param.path = "/v1/me/getexecutions"
	param.method = http.MethodGet
	param.isPrivate = true
	queries := url.Values{}
	queries.Add("product_code", string(client.productCode))
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
