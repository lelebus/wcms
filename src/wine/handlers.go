package wine

import (
	"database/sql"
	"encoding/json"
	"strings"
	"log"
	"net/http"
	"errors"
	"strconv"
	"time"
)

var DB *sql.DB

// Multiplexer for handling /wine requests
func WineHandler (w http.ResponseWriter, r *http.Request) {

	// check correctness of request
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "", 415)
		log.Println(`ERROR in request-header "Content-Type" field: just "application/json" is accepted`)
		return
	}

	switch r.Method {
	case "GET": getWine(w, r)
	case "POST": createWine(w, r) 
	case "PATCH": updateWine(w, r)
	case "DELETE": deleteWine(w, r)
	}
}

// JSON response for request error
func writeError(id string, message string, w http.ResponseWriter) {
	var body = `{ "id":"` + id + `", "message":"` + message + `" }`
	http.Error(w, body, 422)
}

//////////////////////////////////////////////////////////
//
// Handle GET method for wine
//
//////////////////////////////////////////////////////////
func getWine(w http.ResponseWriter, r *http.Request) {
	selection := r.URL.Path[len("/wine/"):]
	
	var query = `SELECT ID, area, type, size, name, winery, year, region, country, price, catalog, details, internalnotes FROM wine `
	
	if selection != "" {
		query += "WHERE id = " + selection
	}
	body, err := queryWine(query)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type","application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func queryWine(query string) ([]byte, error) {

	// query database
	rows, err := DB.Query(query)
	if err != nil {
		err = errors.New("ERROR in retrieving wine entries from DB: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	// read retrieved lines
	wines := make([]Wine, 0)
	for rows.Next() {
		wine := Wine{}

		// FINISH WHEN DB READY
		err = rows.Scan(&wine.ID, &wine.Area, &wine.Type )
		if err != nil {
			err = errors.New("ERROR in scanning retrieved wine entries: " + err.Error())
			return nil, err
		}
		wines = append(wines, wine)
	}

	// marshal wines
	body, err := json.Marshal(wines)
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

	err, wines := readWine(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for _, wine := range wines {  
		id_error, err := checkWine(wine)
		if err != nil {
			writeError(id_error, err.Error(), w)
			log.Println("ERROR in parameter checking: " + err.Error())
			return
		}
	}
	
	for _, wine := range wines {
		insertWine(wine)

		w.WriteHeader(http.StatusOK)
		log.Printf("SUCCESSFUL import: \"%v BY %v - %v\" at line %v \n", wine.Name, wine.Winery, wine.Year, wine.ID)
	}
}

// Create array of Wine from json array given as input
func readWine(r *http.Request) (error, []Wine) {
	var wines []Wine

	decoder := json.NewDecoder(r.Body)

	// read open bracket
	_, err := decoder.Token()
	if err != nil {
		return err, nil
	}

	for decoder.More() {
		var wine Wine
		// decode line 
		err := decoder.Decode(&wine)
		if err != nil {
			return err, nil
		}

		wines = append(wines, wine)
		log.Printf("SUCCESSFUL reading from import JSON:  \"%v BY %v - %v\" \n", wine.Name, wine.Winery, wine.Year)
	}

	// read closing bracket
	_, err = decoder.Token()
	if err != nil {
		return err, nil
	}

	return nil, wines
}

// Check that all parameters of a wine are accepted
func checkWine(wine Wine) (string, error) { 

	if !contains(WineType, strings.ToLower(wine.Type)) {
		e := wine.Type + " is not an accepted TYPE for wine. Check line " + wine.ID
		return "type", errors.New(e)
	} 

	if !contains(WineSize, wine.Size) {
		e := wine.Size + " is not an accepted SIZE for wine (Use . as decimal separator). Check line " + wine.ID
		return "size", errors.New(e)
	}

	dt := time.Now()
	today := dt.Format("02-01-2006")
	currentYear, _ := strconv.ParseInt(today[6:], 10, 64)

	productionYear, err := strconv.ParseInt(wine.Year, 10, 64)
	if err != nil {
		e := "YEAR of wine must be an integer. Check line " + wine.ID
		return "production_year", errors.New(e)
	}
	if  productionYear > currentYear {
		e := "YEAR of wine cannot be set in the future. Check line " + wine.ID
		return "production_year", errors.New(e)
	} 

	v, err := strconv.ParseFloat(wine.Price,10)
	if err != nil {
		e := wine.Price + " is not an accepted PRICE for wine (Must have . as decimal separator). Check line " + wine.ID
		return "price", errors.New(e)
	}
	if v < 0 {
		err := wine.Price + " is not an accepted PRICE for wine (Must be positive). Check line " + wine.ID
		return "price", errors.New(err)
	}

	log.Printf("SUCCESSFUL parameter checking: \"%v BY %v - %v\" at line %v \n", wine.Name, wine.Winery, wine.Year, wine.ID)
	return "", nil
}
// Insert wine in database, checking insertion in other catalogs
func insertWine(wine Wine) {
	// IF ALREADY IN DB, UPDATE WITH INPUT VALUES
	// CHECK IF IT SATISFIES REQUIREMENTS FOR SOME CATALOG
}

//////////////////////////////////////////////////////////
//
// Handle PATCH method for wine update
//
//////////////////////////////////////////////////////////
func updateWine(w http.ResponseWriter, r *http.Request) {
	err, wines := readWine(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	wine := wines[0]

	// UPDATE query

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL update: \"%v BY %v - %v\" \n", wine.Name, wine.Winery, wine.Year)
}

//////////////////////////////////////////////////////////
//
// Handle DELETE method for wine delete
//
//////////////////////////////////////////////////////////
func deleteWine(w http.ResponseWriter, r *http.Request) {
	err, wines := readWine(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	wine := wines[0]

	// DELETE query

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL delete: \"%v BY %v - %v\" \n", wine.Name, wine.Winery, wine.Year)
}