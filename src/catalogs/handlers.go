package catalogs

import (
	// server "WCMS/src/main"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lib/pq"
)

// i = server.DB

var DB *sql.DB

// URLPath for this API
var URLPath = "/catalogs/"

// Multiplexer for handling /catalog requests
func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQUEST Path: %v - Method: %v \n", r.URL.Path, r.Method)

	switch r.Method {
	case "GET":
		getCatalog(w, r)
	case "POST":
		createCatalog(w, r)
	case "PATCH":
		updateCatalog(w, r)
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

	selection := r.URL.Path[len(URLPath):]

	var query string
	var err error
	var body []byte

	if selection == "" {
		query = `SELECT name, level FROM catalog` //QUERY ALL CATALOGS
		body, err = queryCatalog(true, query)
	} else {
		query = `SELECT * WHERE name = ` + selection //QUERY SELECTION
		body, err = queryCatalog(false, query)
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

	/*
		// MOCK UP
		var catalogs []Catalog

		one := Catalog{1, "Stellar Wines", 0, "", []string{"sparkling"}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{"1"}}
		two := Catalog{2, "Vini Italiani", 0, "", []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{"2"}}
		three := Catalog{3, "Vini Italiani / Friuli Venezia Giulia", 1, "Vini Italiani", []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{"2"}}
		if all {
			catalogs = []Catalog{one, two, three}
		} else {
			selection := query[(len(query) - 1):]

			switch selection {
			case "1":
				catalogs = []Catalog{one}
			case "2":
				catalogs = []Catalog{two}
			case "3":
				catalogs = []Catalog{three}
			}
		}
		// END
	*/

	//query database
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
			err = rows.Scan(&catalog.ID, &catalog.Name, &catalog.Level)
		} else {
			err = rows.Scan(&catalog.ID, &catalog.Name, &catalog.Level, &catalog.Parent, pq.Array(&catalog.Type), pq.Array(&catalog.Size), pq.Array(&catalog.Year), pq.Array(&catalog.Territory), pq.Array(&catalog.Region), pq.Array(&catalog.Country), pq.Array(&catalog.Winery), pq.Array(&catalog.Storage), pq.Array(&catalog.Wine))
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
func createCatalog(w http.ResponseWriter, r *http.Request) {

	// check correctness of request
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "", 415)
		log.Println(`ERROR in request-header "Content-Type" field: just "application/json" is accepted`)
		return
	}

	catalogs, err := readCatalogFromJSON(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
	}

	for _, catalog := range catalogs {
		var query string

		if !catalog.Customized {
			// get wines matching catalog parameters
			query = `SELECT id FROM wine WHERE`
			wines, err := getMatchingIDs(query)
			if err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				log.Println(err)
				return
			}
			catalog.Wine = wines
		}

		// insert catalog
		query = `INSERT INTO catalog ( LIST VALUES ) VALUES ($1)`
		_, err = DB.Exec(query, pq.Array(catalog.Name))
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Println("ERROR inserting catalog \"" + catalog.Name + "\"in DB: " + err.Error())
		}

		log.Printf("SUCCESSFUL import: \"%v\"", catalog.Name)
	}

	w.WriteHeader(http.StatusOK)
}

func readCatalogFromJSON(r *http.Request) ([]Catalog, error) {
	var catalogs []Catalog

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e := "ERROR in reading input JSON: " + err.Error()
		return nil, errors.New(e)
	}

	err = json.Unmarshal(body, &catalogs)
	if err != nil {
		e := "ERROR in unmarshalling JSON body: " + err.Error()
		return nil, errors.New(e)
	}

	return catalogs, nil
}

func getMatchingIDs(query string) ([]int, error) {

	//query database
	rows, err := DB.Query(query)
	if err != nil {
		err = errors.New("ERROR in retrieving catalog entries from DB: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	// read retrieved lines
	array := make([]int, 0)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			err = errors.New("ERROR in scanning retrieved catalog entries: " + err.Error())
			return nil, err
		}

		array = append(array, id)
	}
	return array, nil
}

//////////////////////////////////////////////////////////
//
// Handle PATCH method for catalog update
//
//////////////////////////////////////////////////////////
func updateCatalog(w http.ResponseWriter, r *http.Request) {

	// check correctness of request
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "", 415)
		log.Println(`ERROR in request-header "Content-Type" field: just "application/json" is accepted`)
		return
	}

	catalogs, err := readCatalogFromJSON(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for _, catalog := range catalogs {
		// update name
		query := `UPDATE catalog SET name = $1 WHERE id = $2;`
		_, err = DB.Exec(query, catalog.Name, catalog.ID)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Println("ERROR updating catalog name to \"" + catalog.Name + "\"in DB: " + err.Error())
		}

		log.Printf("SUCCESSFUL update: \"%v\" \n", catalog.Name)
	}

	w.WriteHeader(http.StatusOK)
}

//////////////////////////////////////////////////////////
//
// Handle DELETE method for catalog deletion
//
////////////////////////////////////////////////////////
func deleteCatalog(w http.ResponseWriter, r *http.Request) {
	selection := r.URL.Path[len(URLPath):]

	var query string
	var err error

	// ATOMIC FUNCTION to-do

	// delete catalog
	query = `DELETE FROM catalogs WHERE id = $1;`
	_, err = DB.Exec(query, selection)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println("ERROR deleting catalog \"" + selection + "\" from DB: " + err.Error())
	}

	// delete catalog references in wines
	query = `UPDATE FROM catalogs WHERE id = $1;`
	_, err = DB.Exec(query, selection)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println("ERROR deleting catalog \"" + selection + "\" references from DB: " + err.Error())
	}

	// END

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL delete ID: %v \n", selection)
}
