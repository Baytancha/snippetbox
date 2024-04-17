package main

import (
	"errors"
	"fmt"

	//"html/template" // New import
	// New import
	"net/http"
	"strconv"

	"github.com/Baytancha/snip56/internal/models"
	"github.com/julienschmidt/httprouter" // New import
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Because httprouter matches the "/" path exactly, we can now remove the
	// manual check of r.URL.Path != "/" from this handler.

	// if r.URL.Path != "/" { //restricting the wildcard pattern
	// 	app.notFound(w) // Use the notFound() helper
	// 	//http.NotFound(w, r)
	// 	return
	// }

	//panic("oops! something went wrong") // Deliberate panic

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.

	// data := &templateData{
	// 	Snippets: snippets,
	// }

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.page.tmpl", data)

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%+v\n", snippet)
	// }

	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.

	//когда используем Execute(), нужно пользоваться вложенными шаблонами
	//(home реализует шаблон base который реализует шаблоны из home)
	// files := []string{
	// 	"C:\\Users\\mk\\snippetbox\\ui\\html\\home.page.tmpl",
	// 	"C:\\Users\\mk\\snippetbox\\ui\\html\\base.layout.tmpl",
	// 	"C:\\Users\\mk\\snippetbox\\ui\\html\\footer.partial.tmpl",
	// }

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file p
	// as a variadic parameter?
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// Because the home handler function is now a method against application
	// 	// it can access its fields, including the error logger. We'll write the log
	// 	// message to this instead of the standard logger.
	// 	app.serverError(w, err) // Use the serverError() helper.
	// 	//app.errorLog.Println(err.Error())
	// 	//log.Println(err.Error())
	// 	//http.Error(w, "Internal Server Error", 500)
	// 	return
	// }

	// Create an instance of a templateData struct holding the slice of
	// snippets.
	// data := &templateData{
	// 	Snippets: snippets,
	// }

	// // Pass in the templateData struct when executing the template.
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

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
	//err = ts.Execute(w, nil)
	//if err != nil {
	// Also update the code here to use the error logger from the application
	// struct.
	//app.serverError(w, err) // Use the serverError() helper.
	//app.errorLog.Println(err.Error())
	//log.Println(err.Error())
	//http.Error(w, "Internal Server Error", 500)
	//}

	//w.Write([]byte("Hello from Snippetbox"))
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	// When httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context. We'll talk about request context
	// in detail later in the book, but for now it's enough to know that you can
	// use the ParamsFromContext() function to retrieve a slice containing these
	// parameter names and values like so:
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" named
	// parameter from the slice and validate it as normal.
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// if err != nil || id < 1 {
	// 	app.notFound(w) // Use the notFound() helper.
	// 	//http.NotFound(w, r)
	// 	return
	// }

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Initialize a slice containing the paths to the view.tmpl file,
	// plus the base layout and navigation partial that we made earlier.
	// files := []string{
	// 	//"C:\\Users\\mk\\snippetbox\\ui\\html\\view.templ.tmpl", //template instantiation file must be first
	// 	//"C:\\Users\\mk\\snippetbox\\ui\\html\\home.page.tmpl",
	// 	//"C:\\Users\\mk\\snippetbox\\ui\\html\\view.tmpl",
	// 	"C:\\Users\\mk\\snippetbox\\ui\\html\\base.layout.tmpl",
	// 	"C:\\Users\\mk\\snippetbox\\ui\\html\\footer.partial.tmpl",
	// 	"C:\\Users\\mk\\snippetbox\\ui\\html\\view.tmpl", //works with view.templ too
	// 	"C:\\Users\\mk\\snippetbox\\ui\\html\\nav.tmpl",
	// }

	// Parse the template files...
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// Create an instance of a templateData struct holding the snippet data.
	// data := &templateData{
	// 	Snippet: snippet,
	// }

	// And do the same thing again here...
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)

	//при использовании ExecuteTemplate не нужно собирать вложенные шаблоны и соблюдать порядок вызова шаблонов

	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

	// Write the snippet data as a plain-text HTTP response body.
	//fmt.Fprintf(w, "%+v", snippet)

	//fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a new snippetCreate handler, which for now returns a placeholder
// response. We'll update this shortly to show a HTML form.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {

	// Checking if the request method is a POST is now superfluous and can be
	// removed, because this is done automatically by httprouter.
	// if r.Method != "POST" {
	// 	w.Header().Set("Allow", "POST")
	// 	app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
	// 	//http.Error(w, "Method Not Allowed", 405)
	// 	return
	// }
	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	//Using Sprintf for fast dirty concatenations
	//http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

	// Update the redirect path to use the new clean URL format.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

	//w.Write([]byte("Create a new snippet..."))
}
