package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/duybnguyen/snippetbox/pkg/models"
)

// The home handler function.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: snippets}

	// Define the paths to the template files.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Parse the template files.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Execute the template.
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// The showSnippet handler function.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" parameter from the query string and convert it to an integer.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Retrieve the snippet from the database.
	snippet, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Create an instance of a templateData struct holding the snippet data.
	data := &templateData{Snippet: snippet}

	// Define the paths to the template files.
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Parse the template files.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Execute the template.
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// The createSnippet handler function.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Define the snippet data.
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi"
	expires := "7"

	// Insert the snippet into the database.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect to the newly created snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
