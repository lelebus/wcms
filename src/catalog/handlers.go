package catalog

import (
	// server "WCMS/src/main"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// i = server.DB

var DB *sql.DB

func CatalogHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("REQUEST Path: %v - Method: %v \n", r.URL.Path, r.Method)

	// check correctness of request
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "", 415)
		log.Println(`ERROR in request-header "Content-Type" field: just "application/json" is accepted`)
		return
	}

	switch r.Method {
	case "GET":
		getCatalog(w, r)
	case "POST":
		createCatalog(w, r)
	// case "PATCH": updateWine(w, r)
	case "DELETE":
		deleteCatalog(w, r)
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
		query = `SELECT * WHERE name = ` + selection //QUERY SELECTION
	}
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func queryCatalog(all bool, query string) ([]byte, error) {

	// MOCK UP
	selection := query[len(query)-1:]

	var catalogs []Catalog
	one := Catalog{1, "Stellar Wines", 0, "", []string{"sparkling"}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{"1"}}
	two := Catalog{2, "Vini Italiani", 0, "", []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{"2"}}
	three := Catalog{3, "Vini Italiani / Friuli Venezia Giulia", 1, "Vini Italiani", []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{"2"}}
	switch selection {
	case "1":
		catalogs = []Catalog{one}
	case "2":
		catalogs = []Catalog{two}
	case "3":
		catalogs = []Catalog{three}
	default:
		catalogs = []Catalog{one, two, three}
	}
	// END

	// WITH DATABASE
	/*
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
	*/

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
func createCatalog(w http.ResponseWriter, r *http.Request) {
	catalog, err := readCatalog(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
	}

	// GET ALL WINES MATCHING CATALOG PARAMETERS QUERY

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL import: \"%v\"", catalog.Name)
}

func readCatalog(r *http.Request) (Catalog, error) {
	var catalog Catalog

	decoder := json.NewDecoder(r.Body)

	// read open bracket
	_, err := decoder.Token()
	if err != nil {
		return catalog, err
	}

	for decoder.More() {
		// decode line
		err := decoder.Decode(&catalog)
		if err != nil {
			return catalog, err
		}

		log.Printf("SUCCESSFUL reading from import JSON:  \"%v\" \n", catalog.Name)
	}

	// read closing bracket
	_, err = decoder.Token()
	if err != nil {
		return catalog, err
	}

	return catalog, nil
}

//////////////////////////////////////////////////////////
//
// Handle PATCH method for catalog update
//
//////////////////////////////////////////////////////////
func updateCatalog(w http.ResponseWriter, r *http.Request) {
	catalog, err := readCatalog(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// UPDATE query

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL update: \"%v\" \n", catalog.Name)
}

//////////////////////////////////////////////////////////
//
// Handle DELETE method for catalog deletion
//
//////////////////////////////////////////////////////////
func deleteCatalog(w http.ResponseWriter, r *http.Request) {
	selection := r.URL.Path[len("/catalog/"):]

	// DELETE query
	// DELETE FROM WINE

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL delete ID: %v \n", selection)
}
