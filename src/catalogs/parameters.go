package catalogs

import (
	"errors"
	"log"
	"net/http"
)

// URLPath for this API
var ParameterPath = URLPath + "parameters"

// Get all parameter necessary for automatic Catalog creation
func GetAllParameters(w http.ResponseWriter, r *http.Request) {

	log.Printf("REQUEST Path: %v - Method: %v", r.URL.Path, r.Method)

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		log.Println("ERROR in request for \"catalog/parameters/\". Just GET method is allowed")
		return
	}

	body := `{ "parameters": [`

	territories, err := getAll("territory", "origin")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println("ERROR in retrieving Territories: " + err.Error())
		return
	}
	regions, err := getAll("region", "origin")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println("ERROR in retrieving Regions: " + err.Error())
		return
	}
	countries, err := getAll("country", "origin")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println("ERROR in retrieving Countries: " + err.Error())
		return
	}
	wineries, err := getAll("winery", "winery")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println("ERROR in retrieving Wineries: " + err.Error())
		return
	}

	// MOCK UP
	territories = arrayToJSON("territories", []string{"Champagne", "Colli Orientali del Friuli"})
	regions = arrayToJSON("regions", []string{"Friuli Venezia Giulia"})
	countries = arrayToJSON("countries", []string{"France", "Italy"})
	wineries = arrayToJSON("wineries", []string{"Bollinger", "Ronco Severo"})
	// END

	body += wineries + territories + regions + countries + `]}`

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

func getAll(field, table string) (string, error) {
	var fields []string

	query := `SELECT $1 FROM $2;`

	rows, err := DB.Query(query, field, table)
	if err != nil {
		err = errors.New("ERROR in retrieving " + table + " entries from DB: " + err.Error())
		return "", err
	}
	defer rows.Close()

	// read retrieved lines
	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			err = errors.New("ERROR in scanning retrieved + " + table + " entries: " + err.Error())
			return "", err
		}

		fields = append(fields, value)
	}

	return arrayToJSON(field, fields), nil
}
