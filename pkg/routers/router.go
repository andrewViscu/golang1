package routers

import (
	"log"
	"net/http"

	jwt "andrewViscu/golang1/pkg/middleware"

	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter()
	r.Use(setHeaders)
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
	server := &http.Server{Addr: ":1234", Handler: r}
	log.Fatal(server.ListenAndServe())

}

// It's just a test
func setHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)

	})
}
