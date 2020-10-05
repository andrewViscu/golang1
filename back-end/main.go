package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	
)

func main() {
	s := mux.NewRouter()
	

	s.HandleFunc("/getUserProfile", getUserProfile).Methods("GET")
	s.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")
	s.HandleFunc("/createProfile", createProfile).Methods("POST")
	s.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
	
	//TODO: 
	// s.HandleFunc("/", Index).Methods("GET")
	// s.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":1234", s))
}