package routers

import (
	"log"
	"net/http"

	jwt "andrewViscu/golang1/pkg/middleware"

	"github.com/gorilla/mux"
)

//StartServer launches the server
func StartServer() {
	server := &http.Server{Addr: ":1234", Handler: ConfigureServer()}
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
	authRouter.Use(jwt.Handle)

	authRouter.HandleFunc("/{id}", GetUser).Methods("GET")
	authRouter.HandleFunc("/{id}", UpdateUser).Methods("PUT")
	authRouter.HandleFunc("/{id}", DeleteUser).Methods("DELETE")

	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/users", GetAllUsers).Methods("GET")
	r.HandleFunc("/register", CreateUser).Methods("POST")
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/login", TemporaryFunc).Methods("GET")

	return r
}
