package catalogs

import (
	"log"
	"net/http"
)

type origin struct {
	territory string
	region    string
	country   string
}

// URLPath for this API
var ParameterPath = URLPath + "parameters/"

// Get all parameter necessary for automatic Catalog creation
func GetAllParameters(w http.ResponseWriter, r *http.Request) {

	log.Printf("REQUEST Path: %v - Method: %v", r.URL.Path, r.Method)

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		log.Println("ERROR in request for \"catalog/parameters/\". Just GET method is allowed")
		return
	}

	var query string

	// get distinct origins
	var arrayOrigins []origin

	query = `SELECT DISTINCT territory, region, country FROM wine;`

	rows, err := DB.Query(query)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println("ERROR in retrieving origin entries from DB: " + err.Error())
		return
	}
	defer rows.Close()

	// read retrieved lines
	for rows.Next() {
		var row origin
		err = rows.Scan(&row.territory, &row.region, &row.country)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Println("ERROR in scanning retrieved origin entries: " + err.Error())
			return
		}

		arrayOrigins = append(arrayOrigins, row)
	}

	origins := `"origins": [`
	for index, element := range arrayOrigins {
		origins += `{ territory": "` + element.territory + `", "region": "` + element.region + `", "country": "` + element.country + `"}`
		if index != (len(arrayOrigins) - 1) {
			origins += `,`
		}
	}
	origins += `]`

	// get distinct wineries
	var arrayWineries []string

	query = `SELECT DISTINCT winery FROM wine;`

	rows, err = DB.Query(query)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println("ERROR in retrieving winery entries from DB: " + err.Error())
		return
	}
	defer rows.Close()

	// read retrieved lines
	for rows.Next() {
		var row string
		err = rows.Scan(&row)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Println("ERROR in scanning retrieved winery entries: " + err.Error())
			return
		}

		arrayWineries = append(arrayWineries, row)
	}

	wineries := `"wineries": [`
	for index, element := range arrayWineries {
		wineries += `"` + element + `"`
		if index != (len(arrayWineries) - 1) {
			origins += `,`
		}
	}
	wineries += `]`

	body := `{ ` + wineries + `, ` + origins + ` }`

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
