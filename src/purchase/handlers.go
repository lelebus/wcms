package purchase

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Multiplexer for handling /purchase requests
func PurchaseHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("REQUEST Path: %v - Method: %v \n", r.URL.Path, r.Method)

	// check correctness of request
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "", 415)
		log.Println(`ERROR in request-header "Content-Type" field: just "application/json" is accepted`)
		return
	}

	switch r.Method {
	case "GET":
		getPurchase(w, r)
	case "POST":
		createPurchase(w, r)
	case "DELETE":
		deletePurchase(w, r)
	default:
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		log.Println("ERROR in request for \"catalog/parameters/\". Just GET method is allowed")
	}
}

//////////////////////////////////////////////////////////
//
// Handle GET method for purchase
//
//////////////////////////////////////////////////////////
func getPurchase(w http.ResponseWriter, r *http.Request) {
	selection := r.URL.Path[len("/purchase/"):]

	body, err := queryPurchase(selection)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func queryPurchase(id string) ([]byte, error) {

	var purchases []Purchase

	// MOCK UP
	one := Purchase{1, "13/12/1998", "La Tassa", 3, "300.00"}
	purchases = []Purchase{one}
	// END

	// WITH DATABASE
	/*
		//query databases
		rows, err := DB.Query(query)
		if err != nil {
			err = errors.New("ERROR in retrieving purchase entries from DB: " + err.Error())
			return nil, err
		}
		defer rows.Close()

		// read retrieved lines
		purchases := make([]Purchase, 0)
		for rows.Next() {
			purchase := Purchase{}
			err = rows.Scan(&purchase.Date, &purchase.Supplier, &purchase.Quantity, &purchase.Cost)
			if err != nil {
				err = errors.New("ERROR in scanning retrieved purchase entries: " + err.Error())
				return nil, err
			}

			purchases = append(purchases, purchase)
		}
	*/

	// marshal wines
	body, err := json.Marshal(purchases)
	if err != nil {
		err = errors.New("ERROR in marshaling purchase struct to json: " + err.Error())
		return nil, err
	}

	return body, nil
}

//////////////////////////////////////////////////////////
//
// Handle POST method for purchase creation
//
//////////////////////////////////////////////////////////
func createPurchase(w http.ResponseWriter, r *http.Request) {
	purchase, err := readPurchase(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
	}

	idError, err := checkPurchase(purchase)
	if err != nil {
		var body = `{ "id":"` + idError + `", "message":"` + err.Error() + `" }`
		http.Error(w, body, 422)
		log.Println("ERROR in parameter checking: " + err.Error())
		return
	}

	// INSERT PURCHASE QUERY

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL import of purchase for \"%v\" on %v", purchase.Wine, purchase.Date)
}

func readPurchase(r *http.Request) (Purchase, error) {
	var purchase Purchase

	decoder := json.NewDecoder(r.Body)

	// read open bracket
	_, err := decoder.Token()
	if err != nil {
		return purchase, err
	}

	for decoder.More() {
		// decode line
		err := decoder.Decode(&purchase)
		if err != nil {
			return purchase, err
		}

		log.Printf("SUCCESSFUL reading from import JSON: purchase for \"%v\" on %v \n", purchase.Wine, purchase.Date)
	}

	// read closing bracket
	_, err = decoder.Token()
	if err != nil {
		return purchase, err
	}

	return purchase, nil
}

func checkPurchase(purchase Purchase) (string, error) {
	// Date, Quantity, Cost
	dt := time.Now()
	today := dt.Format("02-01-2006")

	if today < purchase.Date {
		e := "DATE of purchase cannot be set in the future."
		return "date", errors.New(e)
	}

	if purchase.Quantity <= 0 {
		e := "QUANTITY of purchase must be and integer."
		return "quantity", errors.New(e)
	}

	price, err := strconv.ParseFloat(purchase.Cost, 10)
	if err != nil {
		e := purchase.Cost + " is not an accepted PRICE for wine (Must have . as decimal separator)."
		return "cost", errors.New(e)
	}

	return "", nil
}
