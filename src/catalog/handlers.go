package catalog

import (
	"net/http"
	"log"
	"database/sql"
	"encoding/json"
	"errors"
)

var DB *sql.DB

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	// check correctness of request
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "", 415)
		log.Println(`ERROR in request-header "Content-Type" field: just "application/json" is accepted`)
		return
	}

	switch r.Method {
	case "GET": getCatalog(w, r)
	// case "POST": createWine(w, r) 
	// case "PATCH": updateWine(w, r)
	// case "DELETE": deleteWine(w, r)
	}
}

//////////////////////////////////////////////////////////
//
// Handle GET method for catalog
//
//////////////////////////////////////////////////////////
func getCatalog(w http.ResponseWriter, r *http.Request) {
	selection := r.URL.Path[len("/wine/"):]
	
	var query string
	var err error
	var body []byte

	if selection == "" {
		query = `SELECT name, level FROM catalog` //QUERY ALL CATALOGS
		body, err = queryCatalog(true, query)
	} else {
		body, err = queryCatalog(false, query)
		query = `SELECT WHERE name = ` + selection  //QUERY SELECTION
	}
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type","application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func queryCatalog(all bool, query string) ([]byte, error) {

	//query databases
	rows, err := DB.Query(query)
	if err != nil {
		err = errors.New("ERROR in retrieving catalog entries from DB: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	// read retrieved lines
	catalogs := make([]Catalog, 0)
	for rows.Next() {
		catalog := Catalog{}

		if all {
			err = rows.Scan(&catalog.Name)
		} else {
			// FINISH WITH DB
			err = rows.Scan(&catalog.Name, &catalog.Type, &catalog.Size)
			//
		}
		if err != nil {
			err = errors.New("ERROR in scanning retrieved catalog entries: " + err.Error())
			return nil, err
		}

		catalogs = append(catalogs, catalog)
	}

	// marshal wines
	body, err := json.Marshal(catalogs)
	if err != nil {
		err = errors.New("ERROR in marshaling catalog struct to json: " + err.Error())
		return nil, err
	}

	return body, nil
}

//////////////////////////////////////////////////////////
//
// Handle POST method for catalog creation
//
//////////////////////////////////////////////////////////
