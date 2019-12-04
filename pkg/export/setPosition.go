package export

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// DB Connection
var DB *sql.DB

// URLPath for this API
var URLPath = "/export/"

type Position struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Position int    `json:"position"`
}

// Multiplexer for handling /wine requests
func SetPosition(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQUEST Path: %v - Method: %v \n", r.URL.Path, r.Method)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		log.Println("ERROR in request for \"catalog/parameters/\". Just GET method is allowed")
	} else {
		// check correctness of request
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "", 415)
			log.Println(`ERROR in request-header "Content-Type" field: just "application/json" is accepted`)
			return
		}

		var position Position

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("ERROR in reading input JSON: " + err.Error())
		}

		err = json.Unmarshal(body, &position)
		if err != nil {
			log.Println("ERROR in unmarshalling JSON body: " + err.Error())
		}

		query := `UPDATE ` + position.Type + ` SET position = $1 WHERE id = $2;`

		_, err = DB.Exec(query, position.Position, position.ID)
		if err != nil {
			log.Println("ERROR setting position for " + position.Type + ": " + string(position.ID) + " :: " + err.Error())
		}

		w.WriteHeader(http.StatusCreated)
		log.Printf("SUCCESSFUL set position for %v: %v", position.Type, position.ID)
	}
}
