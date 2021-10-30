package awin

type DataFeedListRow struct {
	AdvertiserID     string `json:"advertiser_id,omitempty" csv:"Advertiser ID"`
	AdvertiserName   string `json:"advertiser_name,omitempty" csv:"Advertiser Name"`
	PrimaryRegion    string `json:"primary_region,omitempty" csv:"Primary Region"`
	MembershipStatus string `json:"membership_status,omitempty" csv:"Membership Status"`
	FeedID           string `json:"feed_id,omitempty" csv:"Feed ID"`
	FeedName         string `json:"feed_name,omitempty" csv:"Feed Name"`
	Language         string `json:"language,omitempty" csv:"Language"`
	Vertical         string `json:"vertical,omitempty" csv:"Vertical"`
	LastImported     string `json:"last_imported,omitempty" csv:"Last Imported"`
	LastChecked      string `json:"last_checked,omitempty" csv:"Last Checked"`
	NoOfProducts     string `json:"no_of_products,omitempty" csv:"No of products"`
	URL              string `json:"url,omitempty" csv:"URL"`
}
