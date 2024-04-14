package main

import "net/http"

// The routes() method returns a servemux containing our application routes.
// Update the signature for the routes() methods so that it returns a
// http.Handler instead od *http.ServeMux.
func (app *application) routes() http.Handler {
    mux := http.NewServeMux()

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    mux.HandleFunc("/", app.home)
    mux.HandleFunc("/snippet/view", app.snippetView)
    mux.HandleFunc("/snippet/create", app.snippetCreate)

    // Pass the servemux as the 'next' parameter to the secureHeaders middleware.
    // Because secureHeaders is just a function, and the function returns a 
    // http.Handler we don't need to do anything else.
    // wrap the exisiting chain with the logRequest middleware
    return app.logRequest(secureHeaders(mux))
}


// Important: Make sure that you update this signature of the routes() method so that it 
// returns a http.Handler here, otherwise you'll get a compile-time error.
