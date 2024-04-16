package main

import (
	"database/sql"
	"flag"
	"html/template" // New import
	"log"
	"net/http"
	"os"

	//we need the driver’s init() function to run so that it can register itself with the database/sql package.
	_ "github.com/go-sql-driver/mysql"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.01 (Project Setup and Creating
	// a Module) so that the import statement looks like this:
	// "{your-module-path}/internal/models". If you can't remember what module path you
	// used, you can find it at the top of the go.mod file.
	"github.com/Baytancha/snip56/internal/models"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/static/file.zip")
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.

	// Another great feature is that you can use the -help flag to list all the available
	// command-line flags for an application and their accompanying help text. Give it a try:

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	addr := flag.String("addr", "127.0.0.1:4000", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	//handlefunc requires a function wrapped in handler adaptor, but handle requires a handler object

	// Initialize a new instance of our application struct, containing the
	// dependencies.
	app := &application{
		errorLog:      errorLog, //not global vars but accessible via method interfsacing
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	// Call the ListenAndServe() method on our new http.Server struct.
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

	// Write messages using the two new loggers, instead of the standard logger.

	// infoLog.Printf("Starting server on %s", *addr)
	// err := http.ListenAndServe(*addr, mux)
	// errorLog.Fatal(err)

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So we need to dereference the pointer (i.e.
	// prefix it with the * symbol) before using it. Note that we're using the
	// log.Printf() function to interpolate the address with the log message.

	//log.Printf("Starting server on %s", *addr)
	//err := http.ListenAndServe(*addr, mux)
	//log.Fatal(err) //calls os.Exit(1) after writing the message,

	//log.Println("Starting server on :4000")
	//нужно добавлять ip адрес чтобы брандмауэр не жаловался
	//
	//err := http.ListenAndServe("127.0.0.1:4000", mux) //mil instead of mux for the default router
	//log.Fatal(err)
}

// In fact, what exactly is happening is this: When our server receives a new HTTP request,
// it calls the servemux’s ServeHTTP() method. This looks up the relevant handler based on
// the request URL path, and in turn calls that handler’s ServeHTTP() method. You can think
// of a Go web application as a chain of ServeHTTP() methods being called one after another.

// Go's http.Server serves each incoming HTTP request
// in its own goroutine. This allows for efficient handling
// of multiple requests simultaneously without blocking the
// main thread.
// When using middleware functions in Go, such as
// secureHeaders, servemux, and handler, each request
// is processed in its own goroutine to maintain
// concurrency and prevent bottlenecks.

// f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
// if err != nil {
//     log.Fatal(err)
// }
// defer f.Close()

// infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
