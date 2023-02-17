package main

import (
	"fmt"
	"github.com/matthiasbruns/awin-go/awin"
	"net/http"
)

func main() {
	awinClient := awin.NewAwinClient("apiKey", &http.Client{})

	fetchDataFeedList(awinClient)
	fetchDataFeed(awinClient)
}

func fetchDataFeedList(awinClient *awin.AwinClient) {
	feedList, err := awinClient.FetchDataFeedList()

	if err != nil {
		panic(err)
	}

	fmt.Println(feedList)
}

func fetchDataFeed(awinClient *awin.AwinClient) {
	feed, err := awinClient.FetchDataFeed(&awin.DataFeedOptions{
		FeedIds:          []string{"feedId1", "feedId2"},
		Language:         "en",
		ShowAdultContent: false,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(feed)
}
