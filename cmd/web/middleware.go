package main

import (
    "fmt"
    "net/http"
)

// Because we want this middleare to act on every request that is recieved, we need it to be executed before a request hits our servemux. 
// We want the flow of control through our application to look like:
// secureHeaders -> servemux -> application handler
func secureHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Note: This is split accross multiple lines for readability. You don't 
        // need to do this is your own code.
        w.Header().Set("Content-Security-Policy",
            "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

        w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "deny")
        w.Header().Set("X-XSS-Protection", "0")

        next.ServeHTTP(w, r)
    })
}

// Update routes.go file so that logRequest middleware is executed first, and for all requests, 
// so that the flow onf control (reading from left to right) looks like this:
// logRequest <-> secureHeader <-> servemux <-> application handler
func (app *application) logRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

        next.ServeHTTP(w, r)
    })
}

// Create some middleware which recovers the panic and calls our app.serverError() helper method. 
// to do this, we can leverage the fact that deferred functions are always called when the stack is being unwoud following a panic.
func (app *application) recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // create a deffered funciton (which will always be run in the event
        // of a panic as Go unwinds the stack).
        defer func() {
            // Use the buitin recover function to check if there has been a 
            // panic or not. If there was..
            if err := recover(); err != nil {
                // Set a "Connection: close" header on the response.
                w.Header().Set("Connection", "close")
                // Call the app.serverError helper method to return a 500
                // Internal server response.
                app.serverError(w, fmt.Errorf("%s", err))
            }
        }() 
        next.ServeHTTP(w, r)
    })
}




// Setting the Connection: Close header on the response acts as a trigger to make Go's HTTP server automatically close the current connection after a response
// has been sent. It also informs the user that the connection will be closed. Note: if the protocol being used is HTTP/2, Go will automatically strip the 
// Connection: Close header from the response (so it is not malformed) and send a GOAWAY frame.

// The value returned by builtin recover() function has the type any, and its underlying type could be string, error, or something else --whatever the parameter passed to 
// panic() was. In our case, it's the string "opps! something went wrong" (diliberate panic error-- was removed). In the code above we normalized this into an error by using the fmt.Errorf()
// funciton to create a new error object containing the default textual representation of the any value, and the pass this error to the app.serverError() helper method.

