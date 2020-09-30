// package main

// import (
// 	// "github.com/gorilla/mux"
// 	"log"
// 	"net/http"
// 	"github.com/gorilla/mux"
// )

// func main() {
// 	route := mux.NewRouter()
// 	s := route.PathPrefix("/back-end").Subrouter()

// 	s.HandleFunc("/getUserProfile", getUserProfile).Methods("POST")
// 	// s.HandleFunc("/", mainPage).Methods("GET")

// 	//TODO: 
// 	s.HandleFunc("/createProfile", createProfile).Methods("POST")
// 	s.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")
// 	s.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
// 	s.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")
// 	log.Fatal(http.ListenAndServe(":1234", s))
// }