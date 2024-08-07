package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/duybnguyen/snippetbox/pkg/models"
)

// http.ResponseWriter provides methods for assembling a HTTP response and sending it to the user
//http.Request is a struct which holds information about the current request (such as the HTTP method and the URL being requested)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// If we don't return, the handler would keep executing
	if r.URL.Path != "/" {
		app.notFound(w) //implemented in helpers.go
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}
	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"../../ui/html/home.page.tmpl",
		"../../ui/html/base.layout.tmpl",
		"../../ui/html/footer.partial.tmpl",
	}
	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err) //implemented in helpers.go
		return
	}
	// We then use the Execute() method on the template set to write the template
	// content as the response body.
	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err) //implemented in helpers.go
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) //implemented in helpers.go
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string {
		"../../ui/html/show.page.html",
		"../../ui/html/base.layout.tmpl",
		"../../ui/html/footer.partial.tmpl"
	}

	ts.err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err := ts.Execute(w, s)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// must set before WriteHeader and Write
		w.Header().Set("Allow", "POST")
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		app.clientError(w, http.StatusMethodNotAllowed) //implemented in helpers.go
		return
	}

	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}
