package routers

import (
	"log"
	"net/http"

	jwt "andrewViscu/golang1/pkg/middleware"

	"github.com/gorilla/mux"
)

//StartServer launches the server
func StartServer() {
	server := &http.Server{Addr: "1234", Handler: ConfigureServer()}
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func TemporaryFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It's temporary"))
}

func ConfigureServer() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	authRouter := r.PathPrefix("/users").Subrouter()
	restRouter := r.PathPrefix("/").Subrouter()
	restRouter.Headers("Content-Type", "application/json")
	authRouter.Use(jwt.Handle)

	authRouter.HandleFunc("/{id}", GetUser).Methods("GET")
	authRouter.HandleFunc("/{id}", UpdateUser).Methods("PUT")
	authRouter.HandleFunc("/{id}", DeleteUser).Methods("DELETE")

	restRouter.HandleFunc("/", Index).Methods("GET")
	restRouter.HandleFunc("/users", GetAllUsers).Methods("GET")
	restRouter.HandleFunc("/register", CreateUser).Methods("POST")
	restRouter.HandleFunc("/login", Login).Methods("POST")
	restRouter.HandleFunc("/login", TemporaryFunc).Methods("GET")

	return r
}
