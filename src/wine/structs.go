package wine

type Wine struct {
	ID            string `json:"id"`
	Area          string `json:"area"`
	Type          string `json:"type"`
	Size          string `json:"size"`
	Name          string `json:"name"`
	Winery        string `json:"winery"`
	Year          string `json:"year"`
	Territory     string `json:"territory"`
	Region        string `json:"region"`
	Country       string `json:"country"`
	Price         string `json:"price"`
	Catalog       string `json:"catalog"`
	Details       string `json:"details"`
	InternalNotes string `json:"internal-notes"`
	IsActive      bool   `json:"is-active"`
}

// Chiedi a Babbo!!
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
