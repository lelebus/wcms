package catalogs

type Catalog struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Level      int       `json:"level"`
	Parent     int       `json:"parent"`
	Type       []string  `json:"type"`
	Size       []string  `json:"size"`
	Year       []string  `json:"year"`
	Territory  []string  `json:"territory"`
	Region     []string  `json:"region"`
	Country    []string  `json:"country"`
	Winery     []string  `json:"winery"`
	Wines      []string  `json:"wines"`
	Child      []Catalog `json:"child"`
	Customized bool
}
