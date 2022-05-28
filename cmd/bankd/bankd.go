package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lfsgroup/bank"
)

var port = "9200"

func main() {

	p := os.Getenv("PORT")
	if p != "" {
		port = p
	}

	mux := mux.NewRouter()
	mux.HandleFunc("/bsb/{bsb}", LookupBSB)

	mux.Use(JSONMiddleware)

	log.Println("Server starting on port", port)

	bind := ":" + port
	err := http.ListenAndServe(bind, mux)
	if err != nil {
		log.Println("Error starting server: ", err)
	}
}

type errResp struct {
	Message string `json:"message"`
}

// LookupBSB finds a branch via BSB number.
func LookupBSB(wr http.ResponseWriter, req *http.Request) {
	bsb := mux.Vars(req)["bsb"]

	branch, err := bank.LookupBSB(bsb)
	if err != nil {
		switch {
		case errors.Is(err, bank.ErrInvalidBSB):
			wr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(wr).Encode(errResp{
				Message: "invalid BSB number",
			})
		case errors.Is(err, bank.ErrBranchNotFound):
			wr.WriteHeader(http.StatusNotFound)
			json.NewEncoder(wr).Encode(errResp{
				Message: "branch not found",
			})
		default:
			log.Printf("LookupBSB %q error: %v", bsb, err)
			wr.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(wr).Encode(errResp{
				Message: "something went wrong",
			})
		}
		return
	}
	json.NewEncoder(wr).Encode(&branch)
}

// JSON middleware will ensure we only explicitly handle JSON
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		wr.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(wr, req)
	})
}
