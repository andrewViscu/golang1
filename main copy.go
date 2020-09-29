// package main

// import (
// 	"io"
// 	"log"
// 	"net/http"
// )

// func mainPage(w http.ResponseWriter, r *http.Request){
// 	io.WriteString(w, "Hello World\n")
// }

// func main(){
// 	http.HandleFunc("/", mainPage)
// 	log.Fatal(http.ListenAndServe(":1234", nil))
	
// }