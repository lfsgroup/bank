package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/timwmillard/bank"
)

var port = "9200"

func main() {

	p := os.Getenv("PORT")
	if p != "" {
		port = p
	}

	mux := mux.NewRouter()
	mux.HandleFunc("/bsb/{bsb}", LookupBSB)

	log.Println("Server starting on port", port)

	bind := ":" + port
	err := http.ListenAndServe(bind, mux)
	if err != nil {
		log.Println("Error starting server: ", err)
	}
}

// LookupBSB finds a branch via BSB number.
func LookupBSB(wr http.ResponseWriter, req *http.Request) {
	bsb := mux.Vars(req)["bsb"]

	branch, err := bank.LookupBSB(bsb)
	if err != nil {
		wr.WriteHeader(http.StatusNotFound)
	}
	json.NewEncoder(wr).Encode(&branch)
}
