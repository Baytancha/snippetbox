package main

import (
	"fmt"
	"html/template" // New import

	// New import
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { //restricting the wildcard pattern
		app.notFound(w) // Use the notFound() helper
		//http.NotFound(w, r)
		return
	}

	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"C:\\Users\\mk\\snippetbox\\ui\\html\\home.page.tmpl",
		"C:\\Users\\mk\\snippetbox\\ui\\html\\base.layout.tmpl",
		"C:\\Users\\mk\\snippetbox\\ui\\html\\footer.partial.tmpl",
	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file p
	// as a variadic parameter?
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// Because the home handler function is now a method against application
		// it can access its fields, including the error logger. We'll write the log
		// message to this instead of the standard logger.
		app.serverError(w, err) // Use the serverError() helper.
		//app.errorLog.Println(err.Error())
		//log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	// ts, err := template.ParseFiles("C:\\Users\\mk\\snippetbox\\ui\\html\\home.page.tmpl") //different path formatting on windows
	// if err != nil {
	// 	log.Println(err.Error())
	// 	http.Error(w, "Internal Server Error", 500)
	// 	return
	// }

	// We then use the Execute() method on the template set to write the templa
	// content as the response body. The last parameter to Execute() represents
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		// Also update the code here to use the error logger from the application
		// struct.
		app.serverError(w, err) // Use the serverError() helper.
		//app.errorLog.Println(err.Error())
		//log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	w.Write([]byte("Hello from Snippetbox"))
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		//http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		//http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
