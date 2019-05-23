package catalogs

import (
	// server "WCMS/src/main"

	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// URLPath for this API
var GroupedPath = "/catalogs/grouped/"

//////////////////////////////////////////////////////////
//
// Handle GET method for grouped catalog
//
//////////////////////////////////////////////////////////
func GetGroupedCatalogs(w http.ResponseWriter, r *http.Request) {

	var body []byte
	var err error

	query := `SELECT id, name, level, parent, type, size, year, territory, region, country, winery, wines, is_customized FROM catalog WHERE id <> 0;`

	catalogs, err := QueryCatalog(query)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Add children to catalogs
	for _, catalog := range catalogs {
		if catalog.Parent != 0 {
			for pIndex, p := range catalogs {
				if catalog.Parent == p.ID {
					catalogs[pIndex].Child = append(p.Child, catalog)
				}
			}
		}
	}

	// Get top level catalogs
	var topCatalogs []Catalog
	for _, catalog := range catalogs {
		if catalog.Parent == 0 {
			topCatalogs = append(topCatalogs, catalog)
		}
	}

	body, err = json.Marshal(topCatalogs)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		err = errors.New("ERROR in marshaling catalog struct to json: " + err.Error())
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/wcms+json; version=1")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
