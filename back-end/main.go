package main

import (
	// "github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	s := mux.NewRouter()
	

	s.HandleFunc("/getUserProfile", getUserProfile).Methods("GET")
	s.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")
	s.HandleFunc("/createProfile", createProfile).Methods("POST")
	
	//TODO: 
	// s.HandleFunc("/", Index).Methods("GET")
	// s.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
	// s.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":1234", s))
}