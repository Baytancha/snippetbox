package main

import (
	"log"
	"net/http"
)

func main() {
	//handlefunc requires a function, but handle requires a handler object
	mux := http.NewServeMux() //declaring a local servemux
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("C:\\Users\\mk\\snippetbox\\ui\\static"))

	// https://stackoverflow.com/questions/27945310/why-do-i-need-to-use-http-stripprefix-to-access-my-static-files
	// mux.handle передает путь в fileserver. Не просто сопопставляет, а передает
	//Поэтому путь должен быть такой чтобы его понял fileserver, у нас нет второй папки static
	//поэтому удаляем ее из пути и получаем корректный путь к файлу.
	// /static/text.txt -> /text.txt (убираем только /static, второй слэш оставляем!!)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux) //mil instead of mux for the default router
	log.Fatal(err)
}
