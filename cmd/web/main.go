package main

import (
    "flag"
    "log"
    "net/http"
    "os"
)

//Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
    errorLog *log.Logger
    infoLog *log.Logger
}


func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")

    flag.Parse()

    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    
    //Initialize a new instance of our application struct, containing the 
    //dependencies.
    app := &application{
        errorLog: errorLog,
        infoLog: infoLog,
    }

    // Swap the route declaration to use the application struct's methods as the
    // handler functions.
    mux := http.NewServeMux()

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    mux.HandleFunc("/", app.home)
    mux.HandleFunc("/snippet/view", app.snippetView)
    mux.HandleFunc("snippet/create", app.snippetCreate)

    srv := &http.Server{
        Addr: *addr,
        ErrorLog: errorLog,
        Handler: mux,
    }

    infoLog.Printf("Starting serve on %s", *addr)
    err := srv.ListenAndServe()
    errorLog.Fatal(err)
}
