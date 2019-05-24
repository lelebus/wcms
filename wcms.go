package main

import (
	catalog "WCMS/pkg/catalogs"
	purchase "WCMS/pkg/purchases"
	wine "WCMS/pkg/wines"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB
var configuration Config
var port string

func init() {
	configuration = loadConfig("config.json")
	port = ":" + configuration.Port

	// create file for log or append to already existing one
	logfile, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logfile.Close()

	wrt := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(wrt)

	// HOW CAN I SET wrt ALSO FOR THE OTHER PACKAGES??
	// CAN I CREATE DATABASE AUTOMATICALLY FROM HERE

	initDB()
}

func initDB() {
	connection := "host=" + configuration.DB.Host + " port=" + configuration.DB.Port + " dbname=" + configuration.DB.Name +
		" user=" + configuration.DB.User + " password=" + configuration.DB.Password + " sslmode=disable"

	db, err := sql.Open("postgres", connection)
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

	// Initiate DB
	file, err := ioutil.ReadFile("db/init.sql")
	if err != nil {
		err = errors.New("ERROR in reading file for creating DB tables: " + err.Error())
		panic(err)
	}

	_, err = db.Exec(string(file))
	if err != nil {
		err = errors.New("ERROR in initiating postgreSQL database:" + err.Error())
		panic(err)
	}
}

type Config struct {
	Port string `json:"port"`
	DB   struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
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

	return configuration
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

//////////////////////////////////////////////////////////
//
// Start SERVER for REQUEST handling
//
//////////////////////////////////////////////////////////
func main() {
	log.Printf("Starting server at port %v \n", port)

	http.HandleFunc(wine.URLPath, wine.WineHandler)
	http.HandleFunc(purchase.URLPath, purchase.PurchaseHandler)
	http.HandleFunc(catalog.ParameterPath, catalog.GetAllParameters)
	http.HandleFunc(catalog.GroupedPath, catalog.GetGroupedCatalogs)
	http.HandleFunc(catalog.URLPath, catalog.CatalogHandler)
	http.HandleFunc("/", serveJS)

	http.ListenAndServe(port, nil)
}
