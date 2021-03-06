package catalogs

import (
	// server "WCMS/src/main"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

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
	id := r.FormValue("id")
	if id == "" && selection != "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	var body []byte
	var err error

	query := `SELECT id, position, name, level, parent, type, size, territory, region, country, winery, wines, is_customized FROM catalog WHERE `
	if id == "" {
		query += `id <> 0;`
	} else {
		query += `id = ` + id + ";"
	}

	catalogs, err := QueryCatalog(query)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if len(catalogs) == 0 && id != "" {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("ERROR: catalog for given ID can not be found")
		return
	}

	if id != "" {
		body, err = json.Marshal(catalogs[0])
	} else {
		body, err = json.Marshal(catalogs)
	}

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		err = errors.New("ERROR in marshaling catalog struct to json: " + err.Error())
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/wcms+json; version=1")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func QueryCatalog(query string) ([]Catalog, error) {

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

		err = rows.Scan(&catalog.ID, &catalog.Position, &catalog.Name, &catalog.Level, &catalog.Parent, pq.Array(&catalog.Type), pq.Array(&catalog.Size), pq.Array(&catalog.Territory), pq.Array(&catalog.Region), pq.Array(&catalog.Country), pq.Array(&catalog.Winery), pq.Array(&catalog.Wines), &catalog.Customized)
		if err != nil {
			err = errors.New("ERROR in scanning retrieved catalog entries: " + err.Error())
			return nil, err
		}

		catalogs = append(catalogs, catalog)
	}

	return catalogs, nil
}

//////////////////////////////////////////////////////////
//
// Handle POST method for catalog creation
//
//////////////////////////////////////////////////////////
func createCatalog(w http.ResponseWriter, r *http.Request) {

	// check correctness of request
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
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

	ids := []int{}
	for _, catalog := range catalogs {
		if catalog.Level == 0 {
			catalog.Parent = 0
		}
		if catalog.Customized == true {
			catalog.Type = []string{}
			catalog.Size = []string{}
			catalog.Territory = []string{}
			catalog.Region = []string{}
			catalog.Country = []string{}
			catalog.Winery = []string{}
		} else {
			catalog.Wines = []int64{}
		}

		if !catalog.Customized {
			// get wines matching catalog parameters
			wines, err := getMatchingIDs(catalog)
			if err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				log.Println(err)
				return
			}
			catalog.Wines = wines
		}

		log.Println(catalog.Wines)

		id, err := insertCatalog(catalog)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		ids = append(ids, id)

		log.Printf("SUCCESSFUL import: \"%v\"", catalog.Name)
	}

	var body []byte
	if body, err = json.Marshal(ids); err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
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

func getMatchingIDs(catalog Catalog) ([]int64, error) {
	//query database
	query := `
	SELECT DISTINCT w.id FROM wine w WHERE
	( ARRAY[w.type] <@ ($1) OR $1 = '{}' ) AND
	( ARRAY[w.size] <@ ($2) OR $2 = '{}' ) AND
	( ARRAY[w.winery] <@ ($3) OR $3 = '{}' ) AND
	( ARRAY[w.territory] <@ ($4) OR $4 = '{}' ) AND
	( ARRAY[w.region] <@ ($5) OR $5 = '{}' ) AND
	( ARRAY[w.country] <@ ($6) OR $6 = '{}' );`

	rows, err := DB.Query(query, pq.Array(catalog.Type), pq.Array(catalog.Size), pq.Array(catalog.Winery), pq.Array(catalog.Territory), pq.Array(catalog.Region), pq.Array(catalog.Country))
	if err != nil {
		err = errors.New("ERROR in retrieving matching wines to catalog: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	// read retrieved lines
	array := make([]int64, 0)
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			err = errors.New("ERROR in scanning retrieved ids: " + err.Error())
			return nil, err
		}
		array = append(array, id)
	}

	return array, nil
}

func insertCatalog(catalog Catalog) (int, error) {
	var query string
	var err error

	// insert catalog
	query = `INSERT INTO catalog (name, position, level, parent, type, size, territory, region, country, winery, wines, is_customized)
	VALUES ($1, 0, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;`

	result := DB.QueryRow(query, catalog.Name, catalog.Level, catalog.Parent, pq.Array(catalog.Type), pq.Array(catalog.Size), pq.Array(catalog.Territory), pq.Array(catalog.Region), pq.Array(catalog.Country), pq.Array(catalog.Winery), pq.Array(catalog.Wines), catalog.Customized)

	err = result.Scan(&catalog.ID)
	if err != nil {
		err := "ERROR inserting catalog \"" + catalog.Name + "\"in DB: " + err.Error()
		return -1, errors.New(err)
	}

	// insert catalog id in matching wines
	query = `UPDATE wine SET catalogs = array_append(catalogs, $1) WHERE $2 @> ARRAY[id];`

	_, err = DB.Exec(query, catalog.ID, pq.Array(catalog.Wines))
	if err != nil {
		err := "ERROR inserting catalog \"" + catalog.Name + "\" reference in wines: " + err.Error()
		id := strconv.Itoa(catalog.ID)
		deleteFromDB(id)
		return -1, errors.New(err)
	}

	return catalog.ID, nil
}

//////////////////////////////////////////////////////////
//
// Handle PATCH method for catalog update
//
//////////////////////////////////////////////////////////
func updateCatalog(w http.ResponseWriter, r *http.Request) {

	// check correctness of request
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
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

	var query string

	for _, catalog := range catalogs {
		// update name
		if catalog.Customized {
			// Create transaction for DELETE of wine and its references in catalog
			tx, err := DB.Begin()
			if err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				log.Println("ERROR in beginning INSERT procedure for wine" + err.Error())
				return
			}

			query = `UPDATE wine SET catalogs = array_remove(catalogs, $1) WHERE ARRAY[$1]::integer[] <@ catalogs;`
			_, err = DB.Exec(query, catalog.ID)
			if err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				log.Println("ERROR deleting old catalog \"" + catalog.Name + "\" references in wines: " + err.Error())
				return
			}

			query = `UPDATE catalog SET name = $1, wines = $2 WHERE id = $3;`
			_, err = DB.Exec(query, catalog.Name, pq.Array(catalog.Wines), catalog.ID)
			if err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				log.Println("ERROR updating catalog \"" + catalog.Name + "\"in DB: " + err.Error())
				return
			}

			// insert catalog id in matching wines
			query = `UPDATE wine SET catalogs = array_append(catalogs, $1) WHERE $2 @> ARRAY[id];`
			_, err = DB.Exec(query, catalog.ID, pq.Array(catalog.Wines))
			if err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				log.Println("ERROR updating catalog \"" + catalog.Name + "\"in DB: " + err.Error())
				return
			}

			err = tx.Commit()
			if err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				log.Println("ERROR in completing commit for upload PATCH: " + err.Error())
				return
			}

		} else {
			query = `UPDATE catalog SET name = $1 WHERE id = $2;`
			_, err = DB.Exec(query, catalog.Name, catalog.ID)
			if err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				log.Println("ERROR updating catalog \"" + catalog.Name + "\"in DB: " + err.Error())
				return
			}
		}

		log.Printf("SUCCESSFUL update: \"%v\" \n", catalog.Name)
	}

	w.WriteHeader(http.StatusCreated)
}

//////////////////////////////////////////////////////////
//
// Handle DELETE method for catalog deletion
//
////////////////////////////////////////////////////////
func deleteCatalog(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	err := deleteFromDB(id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL delete ID: %v \n", id)
}

func deleteFromDB(id string) error {
	// Begin transaction for DELETE
	tx, err := DB.Begin()
	if err != nil {
		err := "ERROR initializing transaction for catalog DELETE: " + err.Error()
		return errors.New(err)
	}

	var query string

	// delete catalog
	query = `DELETE FROM catalog WHERE id = $1;`

	_, err = DB.Exec(query, id)
	if err != nil {
		err := "ERROR deleting catalog \"" + id + "\" from DB: " + err.Error()
		return errors.New(err)
	}

	// delete catalog references
	query = `UPDATE wine SET catalogs = array_remove(catalogs, $1);`

	_, err = DB.Exec(query, id)
	if err != nil {
		log.Println("ERROR deleting catalog \"" + id + "\" references in wines: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		err := "ERROR in completing commit for catalog DELETE: " + err.Error()
		return errors.New(err)
	}

	return nil
}
