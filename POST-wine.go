package main

import (
	"os"
	"io"
	"encoding/json"
	"strings"
	"log"
	"errors"
	"strconv"
	"time"
)

func init() {
		logfile, err := os.OpenFile("/tmp/wcms.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer logfile.Close()

		wrt := io.MultiWriter(os.Stdout, logfile)
		log.SetOutput(wrt)
}

func main() {

	log.Println("API IMPORT called")

	var input = `
	[		
		{"type": "white", "area": "A 13", "name": "Gaja & Rej", "winery": "Gaja", "year": "2011", "size": "0.75", "region": "langhe - piemonte", "country": "I", "price": "110.00", "catalog": "Speciali Italia", "detail": "", "internal-notes": ""},
		{"type": "red", "area": "A 13", "name": "Gaja Super", "winery": "Gaja", "year": "2015", "size": "0.75", "region": "langhe - piemonte", "country": "I", "price": "250.00", "catalog": "Speciali Italia", "detail": "", "internal-notes": ""}
	]`

	wines := ReadWine(input)

	var linehop int

	switch len(wines) {
	case 0: log.Fatal("ERROR in import: empty file")
	case 1: linehop = 1
	case 2: linehop = 2
	} 

	for line, w := range wines {  
		err := CheckWine(line+linehop, w)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	for line, w := range wines {
		InsertWine(w)

		log.Printf("SUCCESSFUL import: \"%v BY %v - %v\" at line %v \n", w.Name, w.Winery, w.Year, line+linehop)
	}
}

// Create array of Wine from json array given as input
func ReadWine(input string) (wines []Wine) {
	decoder := json.NewDecoder(strings.NewReader(input))

	// read open bracket
	_, err := decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	for decoder.More() {
		var w Wine
		// decode line 
		err := decoder.Decode(&w)
		if err != nil {
			log.Fatal(err)
		}

		wines = append(wines, w)
		log.Printf("SUCCESSFUL reading from import JSON:  \"%v BY %v - %v\" \n", w.Name, w.Winery, w.Year)
	}

	// read closing bracket
	_, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	return
}

// Check that all parameters of a wine are accepted
func CheckWine(line int, w Wine) error { 

	if !contains(WineType, strings.ToLower(w.Type)) {
		err := "ERROR in parameter checking: " + w.Type + " is not an accepted TYPE for wine. Check line " + strconv.Itoa(line)
		return errors.New(err)
	} 

	if !contains(WineSize, w.Size) {
		err := "ERROR in parameter checking: " + w.Size + " is not an accepted SIZE for wine (Use . as decimal separator). Check line " + strconv.Itoa(line)
		return errors.New(err)
	}

	dt := time.Now()
	currentYear := dt.Format("02-01-2006")
	if w.Year > currentYear[6:] {
		err := "ERROR in parameter checking: YEAR of wine cannot be set in the future. Check line " + strconv.Itoa(line)
		return errors.New(err)
	} 

	v, err := strconv.ParseFloat(w.Price,10)
	if err != nil || v < 0 {
		err := "ERROR in parameter checking: " + w.Price + " is not an accepted PRICE for wine (Must be positive and have . as decimal separator. Check line " + strconv.Itoa(line)
		return errors.New(err)
	}

	log.Printf("SUCCESSFUL parameter checking: \"%v BY %v - %v\" at line %v \n", w.Name, w.Winery, w.Year, line)
	return nil
}

// Insert wine in database, checking insertion in other catalogs
func InsertWine(w Wine) {
	
}
