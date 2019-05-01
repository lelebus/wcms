package wine

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DB Connection
var DB *sql.DB

// URLPath for this API
var URLPath string

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
	selection := r.URL.Path[len(URLPath):]

	var query = `SELECT ID, area, type, size, name, winery, year, region, country, price, catalog, details, internalnotes FROM wine `
	var err error
	var body []byte

	// for single wine get more details and purchases
	if selection != "" {
		query += "WHERE id = " + selection
		body, err = queryWine(false, query)
	} else {
		body, err = queryWine(true, query)
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

func queryWine(all bool, query string) ([]byte, error) {

	// MOCK UP
	var wines []Wine
	one := Wine{"1", "X 2", "sparkling", "1.5", "R.D.", "Bollinger", "1985", "Champagne", "", "France", "500", "Stellar Wines", "Recently Disgorged", "", true}
	two := Wine{"2", "Z 14", "white", "0.75", "Ribolla Gialla", "Ronco Severo", "2008", "Colli Orientali del Friuli", "Friuli - Venezia - Giulia", "Italy", "90", "Vini Italiani / Friuli Venezia Giulia", "Macerazione uve", "Alfa Beta Gamma e tu Mamma", true}
	if all {
		wines = []Wine{one, two}
	} else {
		id := query[(len(query) - 1):]
		if id == "1" {
			wines = []Wine{one}
		} else {
			wines = []Wine{two}
		}
	}
	// END

	// WITH DATABASE
	/*
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
				if all {
					err = rows.Scan(&wine.ID, &wine.Area, &wine.Type )
				} else {
					err = rows.Scan(&wine.ID, &wine.Area, &wine.Type )
					// QUERY PURCHASE
				}

			if err != nil {
				err = errors.New("ERROR in scanning retrieved wine entries: " + err.Error())
				return nil, err
			}
			wines = append(wines, wine)
		}
	*/

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
	wines, err := checkWineRequest(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	for _, wine := range wines {
		err := insertWineInDB(wine)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Println(err)
		}

		log.Printf("SUCCESSFUL import: \"%v BY %v - %v\" at line %v \n", wine.Name, wine.Winery, wine.Year, wine.ID)
	}

	w.WriteHeader(http.StatusOK)
}

func checkWineRequest(w http.ResponseWriter, r *http.Request) ([]Wine, error) {

	// check correctness of request
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(415), 415)
		e := `ERROR in request-header "Content-Type" field: just "application/json" is accepted`
		return nil, errors.New(e)
	}

	wines, err := readWineToJSON(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return nil, err
	}

	for _, wine := range wines {
		idError, err := checkWineParameter(wine)
		if err != nil {
			writeError(idError, err.Error(), w)
			e := "ERROR in parameter checking: " + err.Error()
			return nil, errors.New(e)
		}
	}

	return wines, nil
}

// Create array of Wine from json array given as input
func readWineToJSON(r *http.Request) ([]Wine, error) {
	var wines []Wine

	decoder := json.NewDecoder(r.Body)

	// read open bracket
	_, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	for decoder.More() {
		var wine Wine
		// decode line
		err := decoder.Decode(&wine)
		if err != nil {
			return nil, err
		}

		wines = append(wines, wine)
		log.Printf("SUCCESSFUL reading from import JSON:  \"%v BY %v - %v\" \n", wine.Name, wine.Winery, wine.Year)
	}

	// read closing bracket
	_, err = decoder.Token()
	if err != nil {
		return nil, err
	}

	return wines, nil
}

// Check that all parameters of a wine are accepted
func checkWineParameter(wine Wine) (string, error) {

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
	if productionYear > currentYear {
		e := "YEAR of wine cannot be set in the future. Check line " + wine.ID
		return "production_year", errors.New(e)
	}

	v, err := strconv.ParseFloat(wine.Price, 10)
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
func insertWineInDB(wine Wine) error {
	// CHECK IF IT SATISFIES REQUIREMENTS FOR SOME CATALOG
	// INSERT CATALOGS INTO WINE
	// INSERT WINE INTO CATALOGS

	// e := "ERROR in inserting wine \"" + wine.Name + "\" in DB: " + err.Error()
	// return errors.New(e)

	return nil
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
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Printf("ABORTED process for \"%v\" update", wine.Name)
		}

		log.Printf("SUCCESSFUL import: \"%v BY %v - %v\" at line %v \n", wine.Name, wine.Winery, wine.Year, wine.ID)
	}

	w.WriteHeader(http.StatusOK)
}

func updateWineDB(wine Wine) error {
	err := deleteWineFromDB(wine.ID)
	if err != nil {
		return err
	}

	err = insertWineInDB(wine)
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
	selection := r.URL.Path[len(URLPath):]

	err := deleteWineFromDB(selection)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL delete ID: %v \n", selection)
}

func deleteWineFromDB(id string) error {
	// DELETE query
	// DELETE FROM CATALOG

	// e:= "ERROR in ERROR in deleting Wine \""" + wine.Name + "\" from DB: " + err.Error()"
	// return errors.New(e)

	return nil
}
