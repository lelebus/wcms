package main

import (
	wine "WCMS/src/wine"
	"net/http"
	"os"
	"io"
	"log"
)

func init() {
	logfile, err := os.OpenFile("wcms.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logfile.Close()

	wrt := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(wrt)
}

func main() {
	http.HandleFunc("/wine", wine.WineHandler)
	http.ListenAndServe(":8080", nil)
}