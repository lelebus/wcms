package catalogs

import (
	wine "WCMS/pkg/wines"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Parameters struct {
	Types       []string `json:"types"`
	Sizes       []string `json:"sizes"`
	Wineries    []string `json:"wineries"`
	Territories []string `json:"territories"`
	Regions     []string `json:"regions"`
	Countries   []string `json:"countries"`
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

	var p Parameters
	var err error

	p.Types = wine.WineType
	p.Sizes = wine.WineSize

	p.Wineries, err = getParameter("winery", "wine")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	p.Territories, err = getParameter("territory", "wine")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	p.Regions, err = getParameter("region", "wine")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	p.Countries, err = getParameter("country", "wine")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	body, err := json.Marshal(p)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		err := "ERROR in marshaling parameters struct to json: " + err.Error()
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/wcms+json; version=1")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func getParameter(field, table string) ([]string, error) {
	array := []string{}

	query := `SELECT DISTINCT ` + field + ` FROM ` + table + `;`

	rows, err := DB.Query(query)
	if err != nil {
		err := "ERROR in retrieving " + field + " entries from DB: " + err.Error()
		return []string{}, errors.New(err)
	}
	defer rows.Close()

	// read retrieved lines
	for rows.Next() {
		var row string
		err = rows.Scan(&row)
		if err != nil {
			err := "ERROR in scanning retrieved " + field + " entries: " + err.Error()
			return []string{}, errors.New(err)
		}

		array = append(array, row)
	}

	return array, nil
}
