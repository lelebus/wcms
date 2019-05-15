package main

import (
	catalog "WCMS/pkg/catalogs"
	purchase "WCMS/pkg/purchases"
	wine "WCMS/pkg/wines"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB
var port string

func init() {
	logfile, err := os.OpenFile("wcms.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logfile.Close()

	wrt := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(wrt)

	// HOW CAN I SET wrt ALSO FOR THE OTHER PACKAGES??
	// CAN I CREATE DATABASE AUTOMATICALLY FROM HERE

	configuration := loadConfig("config.json")
	port = ":" + configuration.Port

	// connect to database
	connection := "host=" + configuration.DB.Host + " dbname=" + configuration.DB.Name +
		" user=" + configuration.DB.User + " password=" + configuration.DB.Password + " sslmode=disable"

	log.Println(connection)

	db, err = sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	wine.DB = db
	purchase.DB = db
	catalog.DB = db

	log.Println("Successfully connected to database")
}

type Config struct {
	Port string `json:"port"`
	DB   struct {
		Host     string `json:"host"`
		Name     string `json:"name"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"postgreSQL"`
}

func loadConfig(file string) Config {
	var configuration Config

	configFile, err := os.Open(file)
	if err != nil {
		log.Println(err.Error())
		return Config{}
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	decoder.Decode(&configuration)

	// log.Println(configuration)

	return configuration
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
