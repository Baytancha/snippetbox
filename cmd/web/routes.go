package main

import "net/http"

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	//путь
	fileServer := http.FileServer(http.Dir("C:\\Users\\mk\\snippetbox\\ui\\static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	return mux
}
