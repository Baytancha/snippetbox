package main

import "net/http"

// The routes() method returns a servemux containing our application routes.
// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	//путь
	fileServer := http.FileServer(http.Dir("C:\\Users\\mk\\snippetbox\\ui\\static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Wrap the existing chain with the logRequest middleware.
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	//return app.logRequest(secureHeaders(mux))
	//return secureHeaders(mux)
}
