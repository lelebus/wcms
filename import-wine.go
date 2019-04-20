package wine

import (
	"encoding/json"
	"strings"
	"log"
	"net/http"
	"errors"
	"strconv"
	"time"
)

func WineHandler (w http.ResponseWriter, r *http.Request) {
	accept := strings.Split(r.Header.Get("Accept"), ";")
	if accept[0] != "application/json" {
		http.Error(w, "", 415)
		log.Println(`ERROR in Header "Accept" field: just "application/json" is accepted`)
		return
	}

	switch r.Method {
	// case "GET": getWine(r, w)
	default: createWine(w, r) 
	}
}

// Writes JSON response for request error
func writeError(id string, message string, w http.ResponseWriter) {
	var body = `{ "id":"` + id + `", "message":"` + message + `" }`
	http.Error(w, body, 422)
}

// Handles POST method for wine import
func createWine(w http.ResponseWriter, r *http.Request) {
	var input = r.Body()

	// MOCKUP from back-end 
	input = `
	[		
		{"id": "1", "type": "white", "area": "A 13", "name": "Gaja & Rej", "winery": "Gaja", "year": "2011", "size": "0.75", "region": "langhe - piemonte", "country": "I", "price": "110.00", "catalog": "Speciali Italia", "details": "", "internalnotes": ""},
		{"id": "2", "type": "red", "area": "A 13", "name": "Gaja Super", "winery": "Gaja", "year": "2015", "size": "0.75", "region": "langhe - piemonte", "country": "I", "price": "250.00", "catalog": "Speciali Italia", "details": "", "internal-notes": "ALfababab"}
	]`
	// END MOCKUP

	err, wines := readWine(input)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for _, wine := range wines {  
		err = checkWine(wine)
		if err != nil {
			writeError(wine.ID, err.Error(), w)
			log.Println(err)
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
func readWine(input string) (error, []Wine) {
	var wines []Wine

	decoder := json.NewDecoder(strings.NewReader(input))

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
func checkWine(wine Wine) error { 

	if !contains(WineType, strings.ToLower(wine.Type)) {
		err := "ERROR in parameter checking: " + wine.Type + " is not an accepted TYPE for wine. Check line " + wine.ID
		return errors.New(err)
	} 

	if !contains(WineSize, wine.Size) {
		err := "ERROR in parameter checking: " + wine.Size + " is not an accepted SIZE for wine (Use . as decimal separator). Check line " + wine.ID
		return errors.New(err)
	}

	dt := time.Now()
	currentYear := dt.Format("02-01-2006")
	if wine.Year > currentYear[6:] {
		err := "ERROR in parameter checking: YEAR of wine cannot be set in the future. Check line " + wine.ID
		return errors.New(err)
	} 

	v, err := strconv.ParseFloat(wine.Price,10)
	if err != nil || v < 0 {
		err := "ERROR in parameter checking: " + wine.Price + " is not an accepted PRICE for wine (Must be positive and have . as decimal separator. Check line " + wine.ID
		return errors.New(err)
	}

	log.Printf("SUCCESSFUL parameter checking: \"%v BY %v - %v\" at line %v \n", wine.Name, wine.Winery, wine.Year, wine.ID)
	return nil
}

// Insert wine in database, checking insertion in other catalogs
<<<<<<< HEAD:import-file.go
func insertWine(wine Wine) {
	// IF ALREADY IN DB, UPDATE WITH INPUT VALUES
	// CHECK IF IT SATISFIES REQUIREMENTS FOR SOME CATALOG
}
=======
func InsertWine(w Wine) {
	
}
>>>>>>> 8a2bca13f8a00d3d98378d639f5742bc598b5d9e:import-wine.go
