package main

import (
	catalog "WCMS/src/catalogs"
	purchase "WCMS/src/purchases"
	wine "WCMS/src/wines"
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

const port = ":8080"

func init() {
	logfile, err := os.OpenFile("wcms.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logfile.Close()

	wrt := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(wrt)

	// HOW CAN I SET wrt ALSO FOR THE OTHER PACKAGES??
	/*
		// connect to database
		db, err = sql.Open("postgres", "postgres://project:password@localhost/db_project?sslmode=disable")
		if err != nil {
			panic(err)
		}
		if err = db.Ping(); err != nil {
			panic(err)
		}

		wine.DB = db
		purchase.DB = db
		catalog.DB = db
	*/
	log.Println("Successfully connected to database")
}

func main() {
	log.Printf("Starting server at port %v \n", port)

	http.HandleFunc(wine.URLPath, wine.WineHandler)
	http.HandleFunc(purchase.URLPath, purchase.PurchaseHandler)
	http.HandleFunc(catalog.ParameterPath, catalog.GetAllParameters)
	http.HandleFunc(catalog.URLPath, catalog.CatalogHandler)
	http.HandleFunc("/", serveJS)
	http.ListenAndServe(port, nil)
}

func serveJS(w http.ResponseWriter, r *http.Request) {
	filepath := "static/" + r.URL.Path[1:]
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		filepath = "static/index.html"
	}
	http.ServeFile(w, r, filepath)
	log.Printf(`SERVING "%v" for requested path: "%v"`, filepath, r.URL.Path)
}
