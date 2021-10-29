package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/matthiasbruns/awin-go/awin"
	"net/http"
	"os"
	"strings"
)

const cliUsage = "expected 'feedlist' or 'feed' subcommands"
const feedListUsage = "./awin-go feedlist -apikey=API_KEY"
const feedUsage = "./awin-go feed -apikey=API_KEY -ids id1 id2 -lang en -adult true"

func main() {
	awinClient := awin.NewAwinClient(&http.Client{})

	feedListCmd := flag.NewFlagSet("feedlist", flag.ExitOnError)
	feedCmd := flag.NewFlagSet("feed", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println(cliUsage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "feedlist":
		handleFeedListCmd(feedListCmd, awinClient)
	case "feed":
		handleFeedCmd(feedCmd, awinClient)
	default:
		fmt.Println(cliUsage)
		os.Exit(1)
	}
}

func handleFeedListCmd(feedListCmd *flag.FlagSet, awinClient *awin.AwinClient) {
	feedListApiKey := feedListCmd.String("apikey", "", "apikey")

	if err := feedListCmd.Parse(os.Args[2:]); err != nil {
		fmt.Print(feedListUsage)
		os.Exit(1)
	}

	fmt.Println("loading datafeed list from Awin")

	results, err := awinClient.FetchDataFeedList(*feedListApiKey)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	if j, err := json.Marshal(results); err != nil {
		fmt.Print(err)
		os.Exit(1)
	} else {
		fmt.Print(string(j))
	}
}

func handleFeedCmd(feedListCmd *flag.FlagSet, awinClient *awin.AwinClient) {
	feedListApiKey := feedListCmd.String("apikey", "", "-apikey API_KEY")
	feedIds := feedListCmd.String("ids", "", "-ids fleedId1 fleedId2")
	language := feedListCmd.String("lang", "en", "-lang en")
	showAdult := feedListCmd.Bool("adult", false, "-adult true")

	if err := feedListCmd.Parse(os.Args[2:]); err != nil {
		fmt.Print(feedUsage)
		os.Exit(1)
	}

	ids := strings.Split(*feedIds, " ")

	fmt.Println("loading datafeed from Awin")

	results, err := awinClient.FetchDataFeed(&awin.DataFeedOptions{
		ApiKey:           *feedListApiKey,
		FeedIds:          ids,
		Language:         *language,
		ShowAdultContent: *showAdult,
	})

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	if j, err := json.Marshal(results); err != nil {
		fmt.Print(err)
		os.Exit(1)
	} else {
		fmt.Print(string(j))
	}
}
