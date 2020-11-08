package routers

import (
	"log"
	"net/http"

	jwt "andrewViscu/golang1/pkg/middleware"

	"github.com/gorilla/mux"
)
//StartServer launches the server
func StartServer() *http.Server{
	r := mux.NewRouter()

	authRouter := r.PathPrefix("/users").Subrouter()
	restRouter := r.PathPrefix("/").Subrouter()
	restRouter.Headers("Content-Type", "application/json")
	authRouter.Use(jwt.Handle)

	authRouter.HandleFunc("/{id}", GetUser).Methods("GET")
	authRouter.HandleFunc("/{id}", UpdateUser).Methods("PUT")
	authRouter.HandleFunc("/{id}", DeleteUser).Methods("DELETE")

	restRouter.HandleFunc("/",         Index).Methods("GET")
	restRouter.HandleFunc("/users",    GetAllUsers).Methods("GET")
	restRouter.HandleFunc("/register", CreateUser).Methods("POST")
	restRouter.HandleFunc("/login",    Login).Methods("POST")

	server := &http.Server{Addr: ":1234", Handler: r}
	go func() {

        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("ListenAndServe(): %v", err)
        }
    }()

    return server

}
