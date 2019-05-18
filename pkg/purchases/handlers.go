package purchases

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

	//get purchases by wine id
	query := `SELECT p.id, p.wine, p.date, p.supplier, p.quantity, p.cost FROM purchase p, wine w WHERE w.id = $1;`
	rows, err := DB.Query(query, id)
	if err != nil {
		err = errors.New("ERROR in retrieving purchase entries from DB: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	// read retrieved lines
	for rows.Next() {
		purchase := Purchase{}
		var date string

		err = rows.Scan(&purchase.ID, &purchase.Wine, &date, &purchase.Supplier, &purchase.Quantity, &purchase.Cost)
		if err != nil {
			err = errors.New("ERROR in scanning retrieved purchase entries: " + err.Error())
			return nil, err
		}

		// format retrieved date
		t, err := time.Parse("2006-01-02T15:04:05Z", date)
		if err != nil {
			err = errors.New("ERROR in parsing date entry for purchases: " + err.Error())
			return nil, err
		}
		purchase.Date = t.Format("02-01-2006")

		purchases = append(purchases, purchase)
	}

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
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "", 415)
		log.Println(`ERROR in request-header "Content-Type" field: just "application/json" is accepted`)
		return
	}

	purchase, err := readPurchase(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	idError, err := checkPurchase(purchase)
	if err != nil {
		var body = `{ "id":"` + idError + `", "message":"` + err.Error() + `" }`
		http.Error(w, body, 422)
		log.Println("ERROR in parameter checking: " + err.Error())
		return
	}

	err = insertPurchase(purchase)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL import of purchase for \"%v\" on %v", purchase.Wine, purchase.Date)
}

func readPurchase(r *http.Request) (Purchase, error) {
	var purchase Purchase

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e := "ERROR in reading input JSON: " + err.Error()
		return Purchase{}, errors.New(e)
	}

	err = json.Unmarshal(body, &purchase)
	if err != nil {
		e := "ERROR in unmarshalling JSON body: " + err.Error()
		return Purchase{}, errors.New(e)
	}

	return purchase, nil
}

func checkPurchase(purchase Purchase) (string, error) {
	log.Print(purchase)
	dt := time.Now()
	production, err := time.Parse("02-01-2006", purchase.Date)
	if err != nil {
		e := "DATE of purchase cannot be parsed, check request."
		return "date", errors.New(e)
	}

	if dt.Before(production) {
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

func insertPurchase(purchase Purchase) error {
	var query string

	// get highest id inserted for wine
	query = `SELECT max(id) FROM purchase WHERE wine = $1;`

	err := DB.QueryRow(query, purchase.Wine).Scan(&purchase.ID)
	if err != nil {
		purchase.ID = 1
	} else {
		purchase.ID++
	}

	// insert purchase
	_, err = DB.Exec(`SET datestyle to "ISO, DMY";`)

	if err != nil {
		err := "ERROR in inserting purchase: " + err.Error()
		return errors.New(err)
	}

	query = ` INSERT INTO purchase (id, wine, date, supplier, quantity, cost) VALUES ($1, $2, $3, $4, $5, $6);`
	_, err = DB.Exec(query, purchase.ID, purchase.Wine, purchase.Date, purchase.Supplier, purchase.Quantity, purchase.Cost)

	if err != nil {
		err := "ERROR in inserting purchase: " + err.Error()
		return errors.New(err)
	}

	return nil
}

//////////////////////////////////////////////////////////
//
// Handle DELETE method for purchase deletion
//
////////////////////////////////////////////////////////
func deletePurchase(w http.ResponseWriter, r *http.Request) {
	selection := r.URL.Path[len(URLPath):]
	ids := strings.Split(selection, ":")

	// delete purchase
	query := `DELETE FROM purchase WHERE wine = $1 AND id = $2;`
	_, err := DB.Exec(query, ids[0], ids[1])
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		err = errors.New("ERROR in deleting purchase: " + err.Error())
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("SUCCESSFUL delete ID: %v for Wine: %v\n", ids[1], ids[0])
}
