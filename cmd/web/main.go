package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/static/file.zip")
}

func main() {

	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.

	// Another great feature is that you can use the -help flag to list all the available
	// command-line flags for an application and their accompanying help text. Give it a try:

	addr := flag.String("addr", "127.0.0.1:4000", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// 	In staging or production environments, you can redirect the
	// streams to a final destination for viewing and archival.
	// This destination could be on-disk files, or a logging
	// service such as Splunk. Either way, the final destination
	// of the logs can be managed by your execution environment
	// independently of the application.

	// For example, we could redirect the stdout and stderr
	// streams to on-disk files when starting the application
	// like so:
	//go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log

	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flags
	// are joined using the bitwise OR operator |.

	//Custom loggers created by log.New() are concurrency-safe. You can share a single logger
	// and use it across multiple goroutines and in your handlers without needing to worry about race conditions.
	// That said, if you have multiple loggers writing to the same destination then you need to
	// be careful and ensure that the destination’s underlying Write() method is also safe for concurrent use.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//handlefunc requires a function wrapped in handler adaptor, but handle requires a handler object

	// Initialize a new instance of our application struct, containing the
	// dependencies.
	app := &application{
		errorLog: errorLog, //not global vars but accessible via method interfsacing
		infoLog:  infoLog,
	}

	//mux := http.NewServeMux() //declaring a local servemux
	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet", app.showSnippet)
	//mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.

	//fileServer := http.FileServer(http.Dir("C:\\Users\\mk\\snippetbox\\ui\\static"))

	// https://stackoverflow.com/questions/27945310/why-do-i-need-to-use-http-stripprefix-to-access-my-static-files
	// mux передает путь в fileserver. Не просто сопопставляет, а передает
	//Поэтому путь должен быть такой чтобы его понял fileserver, у нас нет второй папки static
	//поэтому удаляем ее из пути и получаем корректный путь к файлу.
	// /static/text.txt -> /text.txt (убираем только /static, второй слэш оставляем!!)
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	//fileserver sanitizes all request paths by running them through the path.Clean() function before searching for a file

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	// Call the ListenAndServe() method on our new http.Server struct.
	err := srv.ListenAndServe()
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
