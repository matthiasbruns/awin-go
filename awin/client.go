// Package awin
// In the awin package you will find all required functions and structs to communicate with the Awin.com services.
package awin

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Constants for url building
const (
	baseUrl = "https://productdata.awin.com"

	/// Example https://productdata.awin.com/datafeed/list/apikey/18a4da1c74680374b05647897c678f94
	dataFeedListUrl = "%s/datafeed/list/apikey/%s"

	/// Example https://productdata.awin.com/datafeed/download/apikey/18a4da1c74680374b05647897c678f94/language/de/fid/123,456/columns/aw_deep_link,product_name,aw_product_id,data_feed_id/format/csv/delimiter/%2C/compression/gzip/adultcontent/1/
	dataFeedUrl = "%s/datafeed/download/apikey/%s/language/%s/fid/%s/columns/%s/format/csv/delimiter/%s/compression/gzip/adultcontent/%d/"
)

var (
	dataFeedColumns = []string{
		"aw_deep_link", "product_name", "aw_product_id", "merchant_product_id", "merchant_image_url", "description",
		"merchant_category", "search_price", "merchant_name", "merchant_id", "category_name", "category_id",
		"aw_image_url", "currency", "store_price", "delivery_cost", "merchant_deep_link", "language", "last_updated",
		"display_price", "data_feed_id", "brand_name", "brand_id", "colour", "product_short_description",
		"specifications", "condition", "product_model", "model_number", "dimensions", "keywords", "promotional_text",
		"product_type", "commission_group", "merchant_product_category_path", "merchant_product_second_category",
		"merchant_product_third_category", "rrp_price", "saving", "savings_percent", "base_price", "base_price_amount",
		"base_price_text", "product_price_old", "delivery_restrictions", "delivery_weight", "warranty",
		"terms_of_contract", "delivery_time", "in_stock", "stock_quantity", "valid_from", "valid_to", "is_for_sale",
		"web_offer", "pre_order", "stock_status", "size_stock_status", "size_stock_amount", "merchant_thumb_url",
		"large_image", "alternate_image", "aw_thumb_url", "alternate_image_two", "alternate_image_three",
		"alternate_image_four", "reviews", "average_rating", "rating", "number_available", "custom_1", "custom_2",
		"custom_3", "custom_4", "custom_5", "custom_6", "custom_7", "custom_8", "custom_9", "ean", "isbn", "upc",
		"mpn", "parent_product_id", "product_GTIN", "basket_link",
	}

	defaultDataFeedColumnsParam string
)

func init() {
	// Generate default column params on init to reduce loops when doing the requests
	defaultDataFeedColumnsParam = strings.Join(dataFeedColumns, ",")
}

// DataFeedOptions
// / FeedIds The string slice of all publisher feed ids
// / Language ISO 3166-1 alpha-2 – two-letter country codes e.g. de, en
// / ShowAdultContent true to include adult content
type DataFeedOptions struct {
	FeedIds          []string
	Language         string
	ShowAdultContent bool
}

// AwinClient
// / client that takes over the communication with the Awin endpoints as well as parsing the response csv data into structs.
// / apiKey You can get the download API key from a standard feed download as given by Create-a-Feed. You can also get the full download link including the relevant API key to access this file from the Create-a-Feed section in the interface (Awin interface --> Toolbox --> Create-a-Feed).
type AwinClient struct {
	client *http.Client
	apiKey string
}

func (c AwinClient) FetchDataFeedList() (*[]DataFeedListRow, error) {
	// Get list of joined and not joined publishers
	resp, err := c.client.Get(fmt.Sprintf(dataFeedListUrl, baseUrl, c.apiKey))
	if err != nil {
		return nil, err
	}

	return parseCSVToDataFeedRow(resp.Body)
}

func (c AwinClient) FetchDataFeed(options *DataFeedOptions) (*[]DataFeedEntry, error) {
	// Get product list of data feed
	showAdult := 0
	if options.ShowAdultContent {
		showAdult = 1
	}

	url := fmt.Sprintf(dataFeedUrl, baseUrl, c.apiKey, options.Language, strings.Join(options.FeedIds, ","), defaultDataFeedColumnsParam, ",", showAdult)

	return c.FetchDataFeedFromUrl(url)
}

func (c AwinClient) FetchDataFeedFromUrl(url string) (*[]DataFeedEntry, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept-Encoding", "gzip")

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	plainResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(plainResponse))
	}

	if err != nil {
		return nil, err
	}

	// gzip response
	gzipReader, err := gzip.NewReader(bytes.NewReader(plainResponse))
	if err != nil {
		return nil, err
	}

	return parseCSVToDataFeedEntry(gzipReader)
}

func parseCSVToDataFeedRow(r io.Reader) (*[]DataFeedListRow, error) {
	var rows []DataFeedListRow

	if err := gocsv.Unmarshal(r, &rows); err != nil {
		return nil, err
	}

	return &rows, nil
}

func parseCSVToDataFeedEntry(r io.Reader) (*[]DataFeedEntry, error) {
	var entries []DataFeedEntry

	if err := gocsv.Unmarshal(r, &entries); err != nil {
		return nil, err
	}

	return &entries, nil
}

func NewAwinClient(apiKey string, client *http.Client) *AwinClient {
	return &AwinClient{client: client, apiKey: apiKey}
}

// NewAwinClientWithHzzp
// / Returns a new NewAwinClientWithHzzp. Needs a http.Client passed from outside.
// / client Required to be passed from the caller
// / returns a new instance of AwinClient
func NewAwinClientWithHttp(apiKey string, client *http.Client) *AwinClient {
	return &AwinClient{client: client, apiKey: apiKey}
}
