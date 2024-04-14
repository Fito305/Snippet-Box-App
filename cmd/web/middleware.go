package main

import (
    "net/http"
)

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


// Because we want this middleare to act on every request that is recieved, we need it to be executed before a request hits our servemux. 
// We want the flow of control through our application to look like:
// secureHeaders -> servemux -> application handler

func (app *application) logRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

        next.ServeHTTP(w, r)
    })
}

// Update routes.go file so that logRequest middleware is executed first, and for all requests, 
// so that the flow onf control (reading from left to right) looks like this:
// logRequest <-> secureHeader <-> servemux <-> application handler
