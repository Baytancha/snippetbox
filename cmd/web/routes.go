package main

import (
	"net/http"

	"github.com/Baytancha/snip56/ui"
	"github.com/julienschmidt/httprouter" // New import
)

// The routes() method returns a servemux containing our application routes.
// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	//mux := http.NewServeMux()
	router := httprouter.New()
	//путь по которому идем в HTML файле
	//fileServer := http.FileServer(http.Dir("C:\\Users\\mk\\snippetbox\\ui\\static"))
	// Take the ui.Files embedded filesystem and convert it to a http.FS type so
	// that it satisfies the http.FileSystem interface. We then pass that to the
	// http.FileServer() function to create the file server handler.
	fileServer := http.FileServer(http.FS(ui.Files))
	//fileserver это корневая папка которая содержит html и static, поэтому нам не надо удалять static из пути
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	//router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Add a new GET /ping route.
	router.HandlerFunc(http.MethodGet, "/ping", ping)
	//гошный роутер не может обрабатывать http-method, он смотрит только на URL, поэтому
	//обработка метода уже осуществляется на уровне хэндлеров
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	//dynamic := alice.New(app.sessionManager.LoadAndSave)

	// And then create the routes using the appropriate methods, patterns and
	// handlers.
	router.Handler(http.MethodGet, "/", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.home))))
	router.Handler(http.MethodGet, "/snippet/view/:id", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.showSnippet)))))
	router.Handler(http.MethodGet, "/user/signup", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.userSignup)))))
	router.Handler(http.MethodPost, "/user/signup", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.userSignupPost)))))
	router.Handler(http.MethodGet, "/user/login", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.userLogin)))))
	router.Handler(http.MethodPost, "/user/login", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.userLoginPost)))))

	//protected := dynamic.Append(app.requireAuthentication)
	//router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))

	router.Handler(http.MethodGet, "/snippet/create", app.sessionManager.LoadAndSave(noSurf(app.authenticate(app.requireAuthentication(http.HandlerFunc(app.createSnippet))))))
	router.Handler(http.MethodPost, "/snippet/create", app.sessionManager.LoadAndSave(noSurf(app.authenticate(app.requireAuthentication(http.HandlerFunc(app.createSnippetPost))))))
	router.Handler(http.MethodPost, "/user/logout", app.sessionManager.LoadAndSave(noSurf(app.authenticate(app.requireAuthentication(http.HandlerFunc(app.userLogoutPost))))))

	//мы попадем на хэндер только если у нас правильный метод
	//router.HandlerFunc(http.MethodGet, "/", app.home)
	//router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.showSnippet)
	//It’s important to be aware that httprouter doesn’t allow conflicting route patterns which potentially
	//match the same request. So, for example, you cannot register a route like GET /foo/new and another
	//route with a named parameter segment or catch-all parameter that conflicts with it —
	//like GET /foo/:name or GET /foo/*name.
	//router.HandlerFunc(http.MethodGet, "/snippet/create", app.createSnippet)
	//router.HandlerFunc(http.MethodPost, "/snippet/create", app.createSnippetPost)

	// Create a handler function which wraps our notFound() helper, and then
	// assign it as the custom handler for 404 Not Found responses. You can also
	// set a custom handler for 405 Method Not Allowed responses by setting
	// router.MethodNotAllowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet/view", app.showSnippet)
	//mux.HandleFunc("/snippet/create", app.createSnippet)

	// Wrap the existing chain with the logRequest middleware.
	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	//return app.logRequest(secureHeaders(mux))
	//return secureHeaders(mux)

	return app.recoverPanic(app.logRequest(secureHeaders(router)))

}
