package main

type Wine struct {
	Area string
	Type string
	Size string
	Name string
	Winery string
	Year string
	Region, Country string
	Price string
	Catalog, Details, InternalNotes string
}

// Chiedi a Babbo!!
var WineType = []string {"red", "white", "sparkling"}
var WineSize = []string {"0.375", "0.75", "1", "1.5", "3", "6", "9", "12"}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
