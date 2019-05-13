package purchases

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

var DB *sql.DB

// URLPath for this API
var URLPath = "/purchases/"

// Multiplexer for handling /purchase requests
func PurchaseHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("REQUEST Path: %v - Method: %v \n", r.URL.Path, r.Method)

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
	selection := r.URL.Path[len(URLPath):]

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

	/*
		// MOCK UP
		one := Purchase{1, 1, "13/12/1998", "La Tassa", 3, "300.00"}
		purchases = []Purchase{one}
		// END
	*/

	//get purchases by wine id
	query := `SELECT p.id, p.date, p.supplier, p.quantity, p.cost FROM purchase p, wine w WHERE w.id = $1;`
	rows, err := DB.Query(query, id)
	if err != nil {
		err = errors.New("ERROR in retrieving purchase entries from DB: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	// read retrieved lines
	for rows.Next() {
		purchase := Purchase{}
		err = rows.Scan(&purchase.ID, &purchase.Date, &purchase.Supplier, &purchase.Quantity, &purchase.Cost)
		if err != nil {
			err = errors.New("ERROR in scanning retrieved purchase entries: " + err.Error())
			return nil, err
		}

		purchases = append(purchases, purchase)
	}

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

	// check correctness of request
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "", 415)
		log.Println(`ERROR in request-header "Content-Type" field: just "application/json" is accepted`)
		return
	}

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

	// insert purchase
	query := `INSERT INTO purchase (id,wine,date,supplier,quantity,cost) VALUES ($1, $2, $3, $4, $5);`
	_, err = DB.Exec(query, purchase.ID, purchase.Wine, purchase.Date, purchase.Supplier, purchase.Quantity, purchase.Cost)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		err = errors.New("ERROR in inserting purchase: " + err.Error())
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL import of purchase for \"%v\" on %v", purchase.Wine, purchase.Date)
}

func readPurchase(r *http.Request) (Purchase, error) {
	var purchase Purchase

	// decoder := json.NewDecoder(r.Body)

	// // read open bracket
	// _, err := decoder.Token()
	// if err != nil {
	// 	return purchase, err
	// }

	// for decoder.More() {
	// 	// decode line
	// 	err := decoder.Decode(&purchase)
	// 	if err != nil {
	// 		return purchase, err
	// 	}

	// 	log.Printf("SUCCESSFUL reading from import JSON: purchase for \"%v\" on %v \n", purchase.Wine, purchase.Date)
	// }

	// // read closing bracket
	// _, err = decoder.Token()
	// if err != nil {
	// 	return purchase, err
	// }

	return purchase, nil
}

func checkPurchase(purchase Purchase) (string, error) {

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
	if price <= 0 {
		e := purchase.Cost + " is not an accepted PRICE for wine (Must be positive)."
		return "cost", errors.New(e)
	}

	return "", nil
}

//////////////////////////////////////////////////////////
//
// Handle DELETE method for purchase deletion
//
////////////////////////////////////////////////////////
func deletePurchase(w http.ResponseWriter, r *http.Request) {
	selection := r.URL.Path[len(URLPath):]
	ids := strings.Split(selection, "-")

	// delete purchase
	query := `DELETE FROM purchase WHERE id = $1 AND wine = $2;`
	_, err := DB.Exec(query, ids[0], ids[1])
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		err = errors.New("ERROR in deleting purchase: " + err.Error())
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL delete ID: %v for Wine: %v\n", ids[0], ids[1])
}
