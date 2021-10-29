package awin

type DataFeedListRow struct {
	AdvertiserID     string `json:"advertiser_id,omitempty"`
	AdvertiserName   string `json:"advertiser_name,omitempty"`
	PrimaryRegion    string `json:"primary_region,omitempty"`
	MembershipStatus string `json:"membership_status,omitempty"`
	FeedID           string `json:"feed_id,omitempty"`
	FeedName         string `json:"feed_name,omitempty"`
	Language         string `json:"language,omitempty"`
	Vertical         string `json:"vertical,omitempty"`
	LastImported     string `json:"last_imported,omitempty"`
	LastChecked      string `json:"last_checked,omitempty"`
	NoOfProducts     string `json:"no_of_products,omitempty"`
	URL              string `json:"url,omitempty"`
}
