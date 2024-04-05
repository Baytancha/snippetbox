package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux() //declaring a local servemux
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux) //mil instead of mux for the default router
	log.Fatal(err)
}
