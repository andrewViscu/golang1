package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	route := mux.NewRouter()
	s := route.PathPrefix("/test").Subrouter()

	s.HandleFunc("/getUserProfile", getUserProfile).Methods("POST")

	//TODO: 
	// s.HandleFunc("/createProfile", createProfile).Methods("POST")
	// s.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")
	// s.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
	// s.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":1234", nil))
}