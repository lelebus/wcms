package catalog

type Catalog struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Level     int
	Parent    string   `json:"parent"`
	Type      []string `json:"type"`
	Size      []string `json:"size"`
	Year      []string `json:"year"`
	Territory []string `json:"territory"`
	Region    []string `json:"region"`
	Country   []string `json:"country"`
	Winery    []string `json:"winery"`
	Storage   []string `json:"storage"`
	Wine      []string `json:"wine"`
}
