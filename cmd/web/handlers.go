package main

import (
	"errors"
	"fmt"

	//"html/template" // New import
	// New import
	"net/http"
	"strconv"

	//"strings"      // New import
	//"unicode/utf8" // New import

	"github.com/Baytancha/snip56/internal/models"
	"github.com/Baytancha/snip56/internal/validator"
	"github.com/julienschmidt/httprouter" // New import
)

// Define a snippetCreateForm struct to represent the form data and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e. start with a capital letter). This is because struct fields
// must be exported in order to be read by the html/template package when
// rendering the template.
type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
	//FieldErrors map[string]string
}

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
	//w.Write([]byte("Display the form for creating a new snippet..."))

	data := app.newTemplateData(r)

	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	// data.Form = snippetCreateForm{
	// 	Expires: 365,
	// }

	app.render(w, http.StatusOK, "create.tmpl", data)

}

// Когда мы закончили заполнять форму мы нажимаем на кнопку и получаем POST-запрос с URL
// Все это зашито в HTML-файл, как только файл генерит URL-запрос, сервер его обрабатывает
// и вызывает функцию createSnippetPost
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

	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError() helper to
	// send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	//title := r.PostForm.Get("title")
	//content := r.PostForm.Get("content")

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an integer. So we need to manually covert
	// the form data to an integer using strconv.Atoi(), and we send a 400 Bad
	// Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create an instance of the snippetCreateForm struct containing the values
	// from the form and an empty map for any validation errors.
	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
		//FieldErrors: map[string]string{},
	}

	// Update the validation checks so that they operate on the snippetCreateForm
	// instance.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// if strings.TrimSpace(form.Title) == "" {
	// 	form.FieldErrors["title"] = "This field cannot be blank"
	// } else if utf8.RuneCountInString(form.Title) > 100 {
	// 	form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	// }

	// if strings.TrimSpace(form.Content) == "" {
	// 	form.FieldErrors["content"] = "This field cannot be blank"
	// }

	// if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
	// 	form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	// }

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "redisplay.tmpl", data)
		return
	}

	// If there are any validation errors re-display the create.tmpl template,
	// passing in the snippetCreateForm instance as dynamic data in the Form
	// field. Note that we use the HTTP status code 422 Unprocessable Entity
	// when sending the response to indicate that there was a validation error.
	// if len(form.FieldErrors) > 0 {
	// 	data := app.newTemplateData(r)
	// 	data.Form = form
	// 	app.render(w, http.StatusUnprocessableEntity, "redisplay.tmpl", data)
	// 	return
	// }

	//ttitle := "O snail"
	//ccontent := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	//eexpires := 7

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	//id, err := app.snippets.Insert(title, content, expires)
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
