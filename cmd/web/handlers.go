package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
)

//Change the signature of the home handler so it is defined as a method against
//*application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        app.notFound(w) // Use the notFound() helper
        return 
    }

    files := []string{
        "./ui/html/base.tmpl.html",
        "./ui/html/partials/nav.tmpl.html",
        "./ui/html/pages/home.tmpl.html",
    }
    
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err) // Use the serverError() helper
        return
    }

    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        app.serverError(w, err) // Use the serverError() helper
    }

}

//Change the signature of the snippetView handler so it is defined as a method
//against *application.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.notFound(w) // USe the app notFound() helper
        return
    }

    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

//Change the signature of the snippetCreate handler so it is defined as a method
//against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
        return
    }

    w.Write([]byte("Create a new snippet..."))
}









// We use the ExecuteTemplate() method to tell Go that we specifically want to respond using the content of hte base template (which in turn invokes our title and main templates.
