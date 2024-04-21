package main

import (
	"net/http"
	
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
// Update the signature for the routes() methods so that it returns a
// http.Handler instead od *http.ServeMux.
func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// Create a handler funciton which wraps our notFOund() helper, and then 
	// assign it as the custom handler for 404 Not Found responses. You can also
	// set a custom handler for 405 Method Not Aloowed responses by setting
	// router.MethodNotALlowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Update the pattern for the route for the static files.
    fileServer := http.FileServer(http.Dir("./ui/static/"))
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// And then create the routes using the appropriate methods, patterns and handlers.
    router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
    router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)
	
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application recieves.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux
	return standard.Then(router)
}


// Important: Make sure that you update this signature of the routes() method so that it 
// returns a http.Handler here, otherwise you'll get a compile-time error.
