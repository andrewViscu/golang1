package routers

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	
)

func StartServer() {
	s := mux.NewRouter()
	

	s.HandleFunc("/users/{id}", GetUser).Methods("GET")
	s.HandleFunc("/users",      GetAllUsers).Methods("GET")
	s.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	// s.HandleFunc("users/{id}",  UpdateUser).Methods("PATCH")
	s.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	s.HandleFunc("/register",   CreateUser).Methods("POST")
	s.HandleFunc("/login",      Login).Methods("POST")

	
	
	//TODO: 
	// s.HandleFunc("/", Index).Methods("GET")
	log.Fatal(http.ListenAndServe(":1234", s))
}
  