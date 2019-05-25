package wines

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

// DB Connection
var DB *sql.DB

// URLPath for this API
var URLPath = "/wines/"

// Multiplexer for handling /wine requests
func WineHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQUEST Path: %v - Method: %v \n", r.URL.Path, r.Method)

	switch r.Method {
	case "GET":
		getWine(w, r)
	case "POST":
		createWine(w, r)
	case "PATCH":
		updateWine(w, r)
	case "DELETE":
		deleteWine(w, r)
	}
}

// JSON response for request error
func writeError(w http.ResponseWriter, reqErr map[string]string) {
	body, err := json.Marshal(reqErr)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		err = errors.New("ERROR in marshaling wine errors to json: " + err.Error())
		log.Println(err)
	}

	http.Error(w, string(body), 422)
}

//////////////////////////////////////////////////////////
//
// Handle GET method for wine
//
//////////////////////////////////////////////////////////
func getWine(w http.ResponseWriter, r *http.Request) {

	selection := r.URL.Path[len(URLPath):]
	id := r.FormValue("id")
	if id == "" && selection != "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	var err error
	var body []byte

	body, err = queryWine(id)
	if err != nil {
		if err.Error() == "404" {
			http.Error(w, http.StatusText(404), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}

		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/wcms+json; version=1")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func queryWine(id string) ([]byte, error) {

	var query = `SELECT id, storage_area, type, size, name, winery, year, territory, region, country, price, catalogs, details, internal_notes FROM wine `
	var rows *sql.Rows
	var err error

	if id != "" {
		query += `WHERE id = $1;`
		rows, err = DB.Query(query, id)
	} else {
		rows, err = DB.Query(query + ";")
	}

	if err != nil {
		return nil, errors.New("404")
	}
	defer rows.Close()

	// read retrieved lines
	wines := make([]Wine, 0)
	for rows.Next() {
		wine := Wine{}

		err = rows.Scan(&wine.ID, &wine.StorageArea, &wine.Type, &wine.Size, &wine.Name, &wine.Winery, &wine.Year, &wine.Territory, &wine.Region, &wine.Country, &wine.Price, pq.Array(&wine.Catalogs), &wine.Details, &wine.InternalNotes)
		if err != nil {
			err = errors.New("ERROR in scanning retrieved wine entries: " + err.Error())
			return nil, err
		}
		if err != nil {
			err = errors.New("ERROR in scanning retrieved wine ID: " + err.Error())
			return nil, err
		}

		wines = append(wines, wine)
	}
	if len(wines) == 0 {
		return nil, errors.New("404")
	}

	// marshal wines
	var body []byte

	if id != "" {
		body, err = json.Marshal(wines[0])
	} else {
		body, err = json.Marshal(wines)
	}

	if err != nil {
		err = errors.New("ERROR in marshaling wine struct to json: " + err.Error())
		return nil, err
	}

	return body, nil
}

//////////////////////////////////////////////////////////
//
// Handle POST method for wine import
//
//////////////////////////////////////////////////////////
func createWine(w http.ResponseWriter, r *http.Request) {
	wines, err := checkWineRequest(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	for _, wine := range wines {
		err := insertWine(wine)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		log.Printf("SUCCESSFUL import: \"%v BY %v - %v\" at line %v \n", wine.Name, wine.Winery, wine.Year, wine.ID)
	}

	w.WriteHeader(http.StatusCreated)
}

func checkWineRequest(w http.ResponseWriter, r *http.Request) ([]Wine, error) {

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, http.StatusText(415), 415)
		e := `ERROR in request-header "Content-Type" field: just "application/json" is accepted`
		return nil, errors.New(e)
	}

	wines, err := readWineFromJSON(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return nil, err
	}

	for _, wine := range wines {
		reqErr := checkWineParameter(wine)
		if len(reqErr) != 0 {
			writeError(w, reqErr)
			e := "ERROR in parameter checking"
			return nil, errors.New(e)
		}
	}

	return wines, nil
}

// Create array of Wine from json array given as input
func readWineFromJSON(r *http.Request) ([]Wine, error) {
	var wines []Wine

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e := "ERROR in reading input JSON: " + err.Error()
		return nil, errors.New(e)
	}

	err = json.Unmarshal(body, &wines)
	if err != nil {
		e := "ERROR in unmarshalling JSON body: " + err.Error()
		return nil, errors.New(e)
	}

	return wines, nil
}

// Check that all parameters of a wine are accepted
func checkWineParameter(wine Wine) map[string]string {

	reqErr := make(map[string]string)

	if wine.StorageArea == "" {
		e := "STORAGE for wine cannot be empty"
		reqErr["storage_area"] = e
	}

	if !contains(WineType, strings.ToLower(wine.Type)) {
		e := wine.Type + ` is not an accepted TYPE for wine`
		reqErr["type"] = e
	}

	if !contains(WineSize, wine.Size) {
		e := wine.Size + ` is not an accepted SIZE for wine (Use . as decimal separator)`
		reqErr["size"] = e
	}

	if wine.Name == "" {
		e := "Wine NAME cannot be empty"
		reqErr["name"] = e
	}

	if wine.Winery == "" {
		e := "Name of WINERY cannot be empty"
		reqErr["winery"] = e
	}

	if wine.Year == "" {
		e := "YEAR of wine must be an integer"
		reqErr["year"] = e
	} else {
		dt := time.Now()
		today := dt.Format("02-01-2006")
		currentYear, _ := strconv.ParseInt(today[6:], 10, 64)

		productionYear, err := strconv.ParseInt(wine.Year, 10, 64)
		if err != nil {
			e := "YEAR of wine must be an integer"
			reqErr["year"] = e
		}
		if productionYear > currentYear {
			e := "YEAR of wine cannot be set in the future"
			reqErr["year"] = e
		}
	}

	if wine.Region == "" && wine.Territory == "" {
		e := "Either REGION or TERRITORY must be set"
		reqErr["region"] = e
		reqErr["territory"] = e
	}

	if wine.Country == "" {
		e := "COUNTRY of wine cannot be empty"
		reqErr["country"] = e
	}

	v, err := strconv.ParseFloat(wine.Price, 10)
	if err != nil {
		e := wine.Price + ` is not an accepted PRICE for wine (Must have . as decimal separator)`
		reqErr["price"] = e
	}
	if v < 0 {
		e := wine.Price + ` is not an accepted PRICE for wine (Must be positive)`
		reqErr["price"] = e
	}

	log.Printf("COMPLETED parameter checking: \"%v BY %v - %v\"\n", wine.Name, wine.Winery, wine.Year)
	return reqErr
}

// Insert wine in database, checking insertion in other catalogs
func insertWine(wine Wine) error {
	log.Println(wine)
	// get catalogs matching wine's parameters
	catalogs, err := getMatchingIDs(wine)
	if err != nil {
		return err
	}

	wine.Catalogs = append(wine.Catalogs, catalogs...)

	var query string

	// insert wine
	if wine.Update {
		query = `
		INSERT INTO wine (id, storage_area,type,size,name,winery,year,territory,region,country,price,catalogs,details,internal_notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`

		_, err = DB.Exec(query, wine.ID, wine.StorageArea, wine.Type, wine.Size, wine.Name, wine.Winery, wine.Year, wine.Territory, wine.Region, wine.Country, wine.Price, pq.Array(wine.Catalogs), wine.Details, wine.InternalNotes)
		if err != nil {
			err := "ERROR inserting wine \"" + wine.Name + "\" in DB: " + err.Error()
			return errors.New(err)
		}
	} else {
		query = `
		INSERT INTO wine (storage_area,type,size,name,winery,year,territory,region,country,price,catalogs,details,internal_notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id;`

		result := DB.QueryRow(query, wine.StorageArea, wine.Type, wine.Size, wine.Name, wine.Winery, wine.Year, wine.Territory, wine.Region, wine.Country, wine.Price, pq.Array(wine.Catalogs), wine.Details, wine.InternalNotes)

		err := result.Scan(&wine.ID)
		if err != nil {
			err := "ERROR inserting wine \"" + wine.Name + "\" in DB: " + err.Error()
			return errors.New(err)
		}

	}

	// insert wine id in matching catalogs
	query = `UPDATE catalog SET wines = array_append(wines, $1) WHERE $2 @> ARRAY[id];`

	_, err = DB.Exec(query, wine.ID, pq.Array(wine.Catalogs))
	if err != nil {
		err := "ERROR inserting wine \"" + wine.Name + "\" in catalogs: " + err.Error()
		id := strconv.Itoa(wine.ID)
		deleteWineFromDB(id)
		return errors.New(err)
	}

	return nil
}

func getMatchingIDs(wine Wine) ([]string, error) {
	//query database
	query := `
	SELECT c.id FROM catalog c WHERE
	( ARRAY[$1] <@ (c.type) OR c.type = '{}' ) AND
	( ARRAY[$2]::float[] <@ (c.size) OR c.size = '{}' ) AND
	( ARRAY[$3]::int[] <@ (c.year) OR c.year = '{}' ) AND
	( ARRAY[$4] <@ (c.territory) OR c.territory = '{}' ) AND
	( ARRAY[$5] <@ (c.region) OR c.region = '{}' ) AND
	( ARRAY[$6] <@ (c.country) OR c.country = '{}' ) AND
	( ARRAY[$7] <@ (c.winery) OR c.winery = '{}' );`

	rows, err := DB.Query(query, wine.Type, wine.Size, wine.Year, wine.Territory, wine.Region, wine.Country, wine.Winery)
	if err != nil {
		err = errors.New("ERROR in retrieving catalog ids matching wine " + wine.Name + " BY " + wine.Winery + " - " + wine.Year + ": " + err.Error())
		return nil, err
	}
	defer rows.Close()

	// read retrieved lines
	array := make([]string, 0)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			err = errors.New("ERROR in scanning retrieved ids: " + err.Error())
			return nil, err
		}

		array = append(array, id)
	}
	return array, nil
}

//////////////////////////////////////////////////////////
//
// Handle PATCH method for wine update
//
//////////////////////////////////////////////////////////
func updateWine(w http.ResponseWriter, r *http.Request) {
	wines, err := checkWineRequest(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	for _, wine := range wines {
		err := updateWineDB(wine)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Printf("ABORTED process for \"%v\" update", wine.Name)
			return
		}

		log.Printf("SUCCESSFUL update: \"%v BY %v - %v\" at line %v \n", wine.Name, wine.Winery, wine.Year, wine.ID)
	}

	w.WriteHeader(http.StatusCreated)
}

func updateWineDB(wine Wine) error {
	// get all customized catalogs, in which wine is inserted
	query := `SELECT id FROM catalog WHERE ARRAY[id] <@ $1 AND is_customized = true;`

	rows, err := DB.Query(query, pq.Array(wine.Catalogs))
	if err != nil {
		err = errors.New("ERROR in retrieving catalog ids matching wine " + string(wine.ID) + ": " + err.Error())
		return err
	}
	defer rows.Close()

	// read retrieved lines
	catalogs := make([]string, 0)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			err = errors.New("ERROR in scanning retrieved catalog ids: " + err.Error())
			return err
		}

		catalogs = append(catalogs, id)
	}
	wine.Catalogs = catalogs
	wine.Update = true

	sel := strconv.Itoa(wine.ID)

	err = deleteWineFromDB(sel)
	if err != nil {
		return err
	}

	err = insertWine(wine)
	if err != nil {
		return err
	}

	return nil
}

//////////////////////////////////////////////////////////
//
// Handle DELETE method for wine delete
//
//////////////////////////////////////////////////////////
func deleteWine(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	err := deleteWineFromDB(id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL delete ID: %v \n", id)
}

func deleteWineFromDB(id string) error {
	// Create transaction for DELETE of wine and its references in catalog
	tx, err := DB.Begin()
	if err != nil {
		err := "ERROR in beginning INSERT procedure for wine" + err.Error()
		return errors.New(err)
	}

	var query string

	// delete wine
	query = `DELETE FROM wine WHERE id = $1;`

	_, err = DB.Exec(query, id)
	if err != nil {
		err := "ERROR deleting wine \"" + id + "\": " + err.Error()
		return errors.New(err)
	}

	// delete wine references
	query = `UPDATE catalog SET wines = array_remove(wines, $1) WHERE ARRAY[$1]::int[] <@ wines`

	_, err = DB.Exec(query, id)
	if err != nil {
		err := "ERROR deleting wine \"" + id + "\" references: " + err.Error()
		return errors.New(err)
	}

	err = tx.Commit()
	if err != nil {
		err := "ERROR in completing commit for wine DELETE: " + err.Error()
		return errors.New(err)
	}

	return nil
}
