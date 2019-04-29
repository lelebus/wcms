package purchase

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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
