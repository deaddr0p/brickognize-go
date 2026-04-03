package brickognize

type Response struct {
	ListingID   string      `json:"listing_id"`
	BoundingBox BoundingBox `json:"bounding_box"`
	Items       []Item      `json:"items"`
}

type BoundingBox struct {
	Left        float64 `json:"left"`
	Upper       float64 `json:"upper"`
	Right       float64 `json:"right"`
	Lower       float64 `json:"lower"`
	ImageWidth  float64 `json:"image_width"`
	ImageHeight float64 `json:"image_height"`
	Score       float64 `json:"score"`
}

type Item struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	ImgURL        string         `json:"img_url"`
	ExternalSites []ExternalSite `json:"external_sites"`
	Category      string         `json:"category"`
	Type          string         `json:"type"`
	Score         float64        `json:"score"`
}

type ExternalSite struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
