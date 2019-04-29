package catalog

import (
	"log"
	"net/http"
)

// Get all parameter necessary for automatic Catalog creation
func GetAllParameter(w http.ResponseWriter, r *http.Request) {

	log.Printf("REQUEST Path: %v - Method: %v", r.URL.Path, r.Method)

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		log.Println("ERROR in request for \"catalog/parameters/\". Just GET method is allowed")
		return
	}

	body := `{ "parameters": [`

	regions, err := getRegions()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	countries, err := getCountries()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	wineries, err := getWineries()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	storage, err := getStorage()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	body += storage + wineries + regions + countries + `]}`

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func arrayToJSON(title string, array []string) string {
	body := `"` + title + `": [`
	for index, element := range array {
		body += `"` + element
		if index != (len(array) - 1) {
			body += `,`
		}
	}
	return body
}

func getRegions() (string, error) {
	var regions []string

	// MOCK UP
	regions = []string{"Champagne", "Colli Orientali - Fiuli Venezia Giulia"}
	// END

	// query all Regions

	return arrayToJSON("regions", regions), nil
}

func getCountries() (string, error) {
	var countries []string

	// MOCK UP
	countries = []string{"France", "Italy"}
	// END

	// query all Countries

	return arrayToJSON("countries", countries), nil
}

func getWineries() (string, error) {
	var wineries []string

	// MOCK UP
	wineries = []string{"Bollinger", "Ronco Severo"}
	// END

	// query all Wineries

	return arrayToJSON("wineries", wineries), nil
}

func getStorage() (string, error) {
	var storage []string

	// MOCK UP
	storage = []string{"X 2", "Z 14"}
	// END

	// query all Wineries

	return arrayToJSON("storage", storage), nil
}
