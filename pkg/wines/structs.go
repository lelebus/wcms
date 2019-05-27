package wines

type Wine struct {
	ID            int    `json:"id"`
	StorageArea   string `json:"storage_area"`
	Type          string `json:"type"`
	Size          string `json:"size"`
	Name          string `json:"name"`
	Winery        string `json:"winery"`
	Year          string `json:"year"`
	Territory     string `json:"territory"`
	Region        string `json:"region"`
	Country       string `json:"country"`
	Price         string `json:"price"`
	Catalogs      []int  `json:"catalog"`
	Details       string `json:"details"`
	InternalNotes string `json:"internal_notes"`
	Update        bool
}

var WineType = []string{"red", "white", "sparkling", "sweet"}
var WineSize = []string{"0.375", "0.75", "1", "1.5", "3", "4.5", "6", "9", "12"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
