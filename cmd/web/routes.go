package main

import (
	"net/http"

	"snippetbox.felipeacosta.net/ui"
	
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

	// Take the ui.Files embedded filesystem and conver it to a https.FS type so that it satisfies the http.FileSystem interface. We then pass that to the http.FileServer() function to create the file server handler.
    fileServer := http.FileServer(http.FS(ui.Files))

	// Our static files are contained in the "static" folder of the ui.Files embedded filesystem. So, for example, our CSS stylesheet is located at "static/css/main.css". This means that we no longer need to strip the prefix from the request URL---any requests that start with /static/ can just be passed directly to the file server and the corresponding static file will be served (so long as it exists).
    router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// Add a new GET /ping route.
	router.HandlerFunc(http.MethodGet, "/ping", ping)
	
	// Create a new midleware chian containing the middleware specific to our dynamic application routes. 
	// For now, this chain will only contain the LoadAndSave session middleware but we'll add more to it.
	// Unprotected application routes usning the 'dynamic' middleware chain.
	// Use the nosurf middleware on all our 'dynamic' routes.
	// Add the authenticate() middleware to the chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// And then create the routes using the appropriate methods, patterns and handlers.
	// Update these routes to use the new dynamic middleware chain followed by the appropriate handler func. Note that becasue the alice ThenFunc() method returns a http.Handler (rather than a http.HanlderFunc) we also need to switch to registering the route using the route.Handler() method.
    router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))


	// Protected (authenticated-only) application status routes, using a new 'protected'
	// middleware chain which includes the requreAuthetication middleware.
	// Because the 'protected' middleware chain appends to the 'dynamic' chain
	// the noSurf middleware will also be used on the three routes below too.
	protected := dynamic.Append(app.requireAuthentication)

    router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application recieves.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux
	return standard.Then(router)
}


// Important: Make sure that you update this signature of the routes() method so that it 
// returns a http.Handler here, otherwise you'll get a compile-time error.
