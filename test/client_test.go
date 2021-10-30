package awin_go

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/matthiasbruns/awin-go/awin"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	response        *http.Response
	requestTestFunc func(r *http.Request) error
}

func (m mockRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if err := m.requestTestFunc(request); err != nil {
		return nil, err
	}
	return m.response, nil
}

func readCSVFileContents(filePath string) (string, error) {
	csvContent, err := ioutil.ReadFile(filePath) // just pass the file name
	if err != nil {
		return "", err
	}
	return string(csvContent), nil
}

func parseCSVToDataFeedRow(csvContent string) (*[]awin.DataFeedListRow, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	var rows []awin.DataFeedListRow
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

		row := awin.DataFeedListRow{
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

func parseCSVToDataFeedEntry(csvContent string) (*[]awin.DataFeedEntry, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	var entries []awin.DataFeedEntry
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

		entry := awin.DataFeedEntry{
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

func TestFetchDataFeedList(t *testing.T) {
	// Read mock data from CSV
	csvContent, err := readCSVFileContents("testdata/data_feed_list.csv")
	if err != nil {
		t.Fatalf("coult not parse csv file '%v'", err)
	}

	// Create mock response
	response := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(csvContent)),
	}

	// Create test client to run tests on
	awinClient := awin.NewAwinClient(&http.Client{Transport: mockRoundTripper{response: response, requestTestFunc: func(r *http.Request) error {
		expectedUrl := "https://productdata.awin.com/datafeed/list/apikey/apiKey"
		if r.URL.String() != expectedUrl {
			err := errors.New(fmt.Sprintf("invalid url found in test\nexpected '%s'\nfound '%s'", expectedUrl, r.URL.String()))
			t.Error(err)
			return err
		}

		expectedMethod := "GET"
		if r.Method != expectedMethod {
			err := errors.New(fmt.Sprintf("invalid request method in test\nexpected '%s'\nfound '%s'", expectedMethod, r.Method))
			t.Error(err)
			return err
		}

		return nil
	}}})

	result, err := awinClient.FetchDataFeedList("apiKey")
	if err != nil {
		t.Fatalf("err is not null '%v'", err)
	}

	if len(*result) != 10 {
		t.Fatalf("Invalid amount of data rows received %d", len(*result))
	}

	// Check if received rows and expected rows match
	expectedRows, _ := parseCSVToDataFeedRow(csvContent)
	for i, expectedRow := range *expectedRows {
		receivedRow := (*result)[i]
		if expectedRow != receivedRow {
			t.Fatalf("Invalid row parsed\nexpected '%v'\nreceived '%v'", expectedRow, receivedRow)
		}
	}
}

func TestFetchDataFeed(t *testing.T) {
	// Read mock data from CSV
	csvContent, err := readCSVFileContents("testdata/data_feed.csv")
	if err != nil {
		t.Fatalf("coult not parse csv file '%v'", err)
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(csvContent)); err != nil {
		t.Error(err)
	}
	if err := gz.Flush(); err != nil {
		t.Error(err)
	}
	if err := gz.Close(); err != nil {
		t.Error(err)
	}

	// Create mock response
	response := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer(b.Bytes())),
	}

	// Create test client to run tests on
	awinClient := awin.NewAwinClient(&http.Client{Transport: mockRoundTripper{response: response, requestTestFunc: func(r *http.Request) error {
		expectedUrl := "https://productdata.awin.com/datafeed/download/apikey/apiKey/language/en/fid/fid1,fid2/columns/aw_deep_link,product_name,aw_product_id,merchant_product_id,merchant_image_url,description,merchant_category,search_price,merchant_name,merchant_id,category_name,category_id,aw_image_url,currency,store_price,delivery_cost,merchant_deep_link,language,last_updated,display_price,data_feed_id,brand_name,brand_id,colour,product_short_description,specifications,condition,product_model,model_number,dimensions,keywords,promotional_text,product_type,commission_group,merchant_product_category_path,merchant_product_second_category,merchant_product_third_category,rrp_price,saving,savings_percent,base_price,base_price_amount,base_price_text,product_price_old,delivery_restrictions,delivery_weight,warranty,terms_of_contract,delivery_time,in_stock,stock_quantity,valid_from,valid_to,is_for_sale,web_offer,pre_order,stock_status,size_stock_status,size_stock_amount,merchant_thumb_url,large_image,alternate_image,aw_thumb_url,alternate_image_two,alternate_image_three,alternate_image_four,reviews,average_rating,rating,number_available,custom_1,custom_2,custom_3,custom_4,custom_5,custom_6,custom_7,custom_8,custom_9,ean,isbn,upc,mpn,parent_product_id,product_GTIN,basket_link/format/csv/delimiter/,/compression/gzip/adultcontent/1/"
		if r.URL.String() != expectedUrl {
			err := errors.New(fmt.Sprintf("invalid url found in test\nexpected '%s'\nfound '%s'", expectedUrl, r.URL.String()))
			t.Error(err)
			return err
		}

		expectedMethod := "GET"
		if r.Method != expectedMethod {
			err := errors.New(fmt.Sprintf("invalid request method in test\nexpected '%s'\nfound '%s'", expectedMethod, r.Method))
			t.Error(err)
			return err
		}

		return nil
	}}})

	result, err := awinClient.FetchDataFeed(&awin.DataFeedOptions{
		ApiKey:           "apiKey",
		FeedIds:          []string{"fid1", "fid2"},
		Language:         "en",
		ShowAdultContent: true,
	})
	if err != nil {
		t.Fatalf("err is not null '%v'", err)
	}

	if len(*result) != 10 {
		t.Fatalf("Invalid amount of data rows received %d", len(*result))
	}

	// Check if received rows and expected rows match
	expectedRows, _ := parseCSVToDataFeedEntry(csvContent)
	for i, expectedRow := range *expectedRows {
		receivedRow := (*result)[i]
		if expectedRow != receivedRow {
			t.Fatalf("Invalid row parsed\nexpected '%v'\nreceived '%v'", expectedRow, receivedRow)
		}
	}
}
