package main

import (
	"html/template" // New import
	"path/filepath" // New import
	"time"          // New import

	"github.com/Baytancha/snip56/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	Snippet     *models.Snippet   //сниппет это связная совокупность данных таблицы
	Snippets    []*models.Snippet //для того чтобы отображать последние n сниппетов
	CurrentYear int
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
// чтобы функция работала в шаблоне она должна возвращать одно значение
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.

// чтобы зарегать функцию в таблице шаблонов нужно засунуть ее в карту
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("C:\\Users\\mk\\snippetbox\\ui\\html\\pages\\*.tmpl") //"./ui/html/pages/*.tmpl"
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles("C:\\Users\\mk\\snippetbox\\ui\\html\\base.layout.tmpl")
		//ts, err := template.ParseFiles("C:\\Users\\mk\\snippetbox\\ui\\html\\base.layout.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("C:\\Users\\mk\\snippetbox\\ui\\html\\partials\\*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the  page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Create a slice containing the filepaths for our base template, any
		// partials and the page.
		// files := []string{
		// 	"C:\\Users\\mk\\snippetbox\\ui\\html\\base.layout.tmpl",
		// 	"C:\\Users\\mk\\snippetbox\\ui\\html\\footer.partial.tmpl",
		// 	"C:\\Users\\mk\\snippetbox\\ui\\html\\nav.tmpl",
		// 	page,
		// }

		// Parse the files into a template set.
		// ts, err := template.ParseFiles(files...)
		// if err != nil {
		// 	return nil, err
		// }

		// Add the template set to the map, using the name of the page
		// (like 'home.tmpl') as the key.
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
