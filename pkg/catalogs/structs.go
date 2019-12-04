package catalogs

type Catalog struct {
	ID         int       `json:"id"`
	Position   int       `json:"position"`
	Name       string    `json:"name"`
	Level      int       `json:"level"`
	Parent     int       `json:"parent"`
	Type       []string  `json:"type"`
	Size       []string  `json:"size"`
	Territory  []string  `json:"territory"`
	Region     []string  `json:"region"`
	Country    []string  `json:"country"`
	Winery     []string  `json:"winery"`
	Wines      []int64   `json:"wines"`
	Child      []Catalog `json:"child"`
	Customized bool      `json:"Customized"`
}
