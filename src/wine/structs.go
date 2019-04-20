package wine

type Wine struct {
	ID string
	Area string
	Type string
	Size string
	Name string
	Winery string
	Year string
	Region, Country string
	Price string
	Catalog, Details string
	InternalNotes string `json:"internal-notes"`
}

// Chiedi a Babbo!!
var WineType = []string {"red", "white", "sparkling", "sweet"}
var WineSize = []string {"0.375", "0.75", "1", "1.5", "3", "4.5", "6", "9", "12"}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
