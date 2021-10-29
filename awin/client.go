// Package awin
// In the awin package you will find all required functions and structs to communicate with the Awin.com services.
package awin

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"errors"
	"fmt"
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
/// ApiKey You can get the download API key from a standard feed download as given by Create-a-Feed. You can also get the full download link including the relevant API key to access this file from the Create-a-Feed section in the interface (Awin interface --> Toolbox --> Create-a-Feed).
/// FeedIds The string slice of all publisher feed ids
/// Language ISO 3166-1 alpha-2 – two-letter country codes e.g. de, en
/// ShowAdultContent true to include adult content
type DataFeedOptions struct {
	ApiKey           string
	FeedIds          []string
	Language         string
	ShowAdultContent bool
}

// AwinClient
/// Client that takes over the communication with the Awin endpoints as well as parsing the response csv data into structs.
type AwinClient struct {
	client *http.Client
}

func (c AwinClient) FetchDataFeedList(apiKey string) (*[]DataFeedListRow, error) {
	// Get list of joined and not joined publishers
	resp, err := c.client.Get(fmt.Sprintf(dataFeedListUrl, baseUrl, apiKey))
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

	url := fmt.Sprintf(dataFeedUrl, baseUrl, options.ApiKey, options.Language, strings.Join(options.FeedIds, ","), defaultDataFeedColumnsParam, ",", showAdult)
	fmt.Println(url)

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
	reader := csv.NewReader(r)
	var rows []DataFeedListRow
	columnNamesSkipped := false
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Skip column names from csv
		if !columnNamesSkipped {
			columnNamesSkipped = true
			continue
		}

		row := DataFeedListRow{
			AdvertiserID:     record[0],
			AdvertiserName:   record[1],
			PrimaryRegion:    record[2],
			MembershipStatus: record[3],
			FeedID:           record[4],
			FeedName:         record[5],
			Language:         record[6],
			Vertical:         record[7],
			LastImported:     record[8],
			LastChecked:      record[9],
			NoOfProducts:     record[10],
			URL:              record[11],
		}

		rows = append(rows, row)
	}

	return &rows, nil
}

func parseCSVToDataFeedEntry(r io.Reader) (*[]DataFeedEntry, error) {
	reader := csv.NewReader(r)
	var entries []DataFeedEntry
	columnNamesSkipped := false
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Skip column names from csv
		if !columnNamesSkipped {
			columnNamesSkipped = true
			continue
		}

		entry := DataFeedEntry{
			AwDeepLink:                    record[0],
			ProductName:                   record[1],
			AwProductId:                   record[2],
			MerchantProductId:             record[3],
			MerchantImageUrl:              record[4],
			Description:                   record[5],
			MerchantCategory:              record[6],
			SearchPrice:                   record[7],
			MerchantName:                  record[8],
			MerchantId:                    record[9],
			CategoryName:                  record[10],
			CategoryId:                    record[11],
			AwImageUrl:                    record[12],
			Currency:                      record[13],
			StorePrice:                    record[14],
			DeliveryCost:                  record[15],
			MerchantDeepLink:              record[16],
			Language:                      record[17],
			LastUpdated:                   record[18],
			DisplayPrice:                  record[19],
			DataFeedId:                    record[20],
			BrandName:                     record[21],
			BrandId:                       record[22],
			Colour:                        record[23],
			ProductShortDescription:       record[24],
			Specifications:                record[25],
			Condition:                     record[26],
			ProductModel:                  record[27],
			ModelNumber:                   record[28],
			Dimensions:                    record[29],
			Keywords:                      record[30],
			PromotionalText:               record[31],
			ProductType:                   record[32],
			CommissionGroup:               record[33],
			MerchantProductCategoryPath:   record[34],
			MerchantProductSecondCategory: record[35],
			MerchantProductThirdCategory:  record[36],
			RrpPrice:                      record[37],
			Saving:                        record[38],
			SavingsPercent:                record[39],
			BasePrice:                     record[40],
			BasePriceAmount:               record[41],
			BasePriceText:                 record[42],
			ProductPriceOld:               record[43],
			DeliveryRestrictions:          record[44],
			DeliveryWeight:                record[45],
			Warranty:                      record[46],
			TermsOfContract:               record[47],
			DeliveryTime:                  record[48],
			InStock:                       record[49],
			StockQuantity:                 record[50],
			ValidFrom:                     record[51],
			ValidTo:                       record[52],
			IsForSale:                     record[53],
			WebOffer:                      record[54],
			PreOrder:                      record[55],
			StockStatus:                   record[56],
			SizeStockStatus:               record[57],
			SizeStockAmount:               record[58],
			MerchantThumbUrl:              record[59],
			LargeImage:                    record[60],
			AlternateImage:                record[61],
			AwThumbUrl:                    record[62],
			AlternateImageTwo:             record[63],
			AlternateImageThree:           record[64],
			AlternateImageFour:            record[65],
			Reviews:                       record[66],
			AverageRating:                 record[67],
			Rating:                        record[68],
			NumberAvailable:               record[69],
			Custom1:                       record[70],
			Custom2:                       record[71],
			Custom3:                       record[72],
			Custom4:                       record[73],
			Custom5:                       record[74],
			Custom6:                       record[75],
			Custom7:                       record[76],
			Custom8:                       record[77],
			Custom9:                       record[78],
			Ean:                           record[79],
			Isbn:                          record[80],
			Upc:                           record[81],
			Mpn:                           record[82],
			ParentProductId:               record[83],
			ProductGtin:                   record[84],
			BasketLink:                    record[85],
		}

		entries = append(entries, entry)
	}

	return &entries, nil
}

// NewAwinClient
/// Returns a new AwinClient. Needs a http.Client passed from outside.
/// client Required to be passed from the caller
/// returns a new instance of AwinClient
func NewAwinClient(client *http.Client) *AwinClient {
	return &AwinClient{client: client}
}
