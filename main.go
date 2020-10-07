package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	
)

func main() {
	s := mux.NewRouter()
	

	s.HandleFunc("/users/{id}", getUserProfile).Methods("GET")
	s.HandleFunc("/users", getAllUsers).Methods("GET")
	s.HandleFunc("/users", createProfile).Methods("POST")
	s.HandleFunc("/users/{id}", updateProfile).Methods("PUT")
	s.HandleFunc("users/{id}", updateProfile).Methods("PATCH")
	s.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")
	
	//TODO: 
	// s.HandleFunc("/", Index).Methods("GET")
	log.Fatal(http.ListenAndServe(":1234", s))
}