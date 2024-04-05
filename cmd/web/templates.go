package main

import (
    "html/template"
    "path/filepath"

    "snippetbox.felipeacosta.net/internal/models"
)


// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
// Include a Snippets field in the templateData struct
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}


func newTemplateCache() (map[string]*template.Template, error) {
    // Initialize a new map to act as the cache.
    cache := map[string]*template.Template{}

    // Use the  filepath.Glob() function to get a slice of all filepaths that
    // match the pattern "./ui/html/pages/*.tmpl. This will essentially give 
    // us a slice of all the filepaths for our applications 'page' templates
    // like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
    pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
    if err != nil {
        return nil, err
    }
    
    for _, page := range pages {
        name := filepath.Base(page)

        // Parse the files into a template set.
        ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
        if err != nil {
            return nil, err
        }

        // Call ParseGlob() *on this template set* to add any partials.
        ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
        if err != nil {
            return nil, err
        }


        // Call ParseFiles() *on this template set* to add the page template.
        ts, err = ts.ParseFiles(page)
        if err != nil {
            return nil, err
        }

        // Add the template set to the map, using the name of the page
        // (like 'home.tmpl') as the key.
        // Add the template set to the map as normal...
        cache[name] = ts
    }

    // Return the map.
    return cache, nil


}
   





    // // Lopp through the page filepaths one-by-one. Old loop
    // for _, page := range pages {
    //     // Extract the file name (like "home.tmpl") from the full filepath
    //     // and assign it to the name varaible.
    //     name := filepath.Base(page)
    //
    //     // Create a slice containing the filepaths for our base template, any 
    //     // partials and the pages.
    //     files := []string{
    //         "./ui/html/base.tmpl.html",
    //         "./ui/html/partials/nav.tmpl.html",
    //         page,
    //     }
