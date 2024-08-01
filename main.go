package main

import (
	"log"
	"net/http"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.

// http.ResponseWriter provides methods for assembling a HTTP response and sending it to the user
//http.Request is a struct which holds information about the current request (such as the HTTP method and the URL being requested)
func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the hand
	// would keep executing and also write the "Hello from SnippetBox" message.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	// If it's not, use the w.WriteHeader() method to send a 405 status code and
	// the w.Write() method to write a "Method Not Allowed" response body. We
	// then return from the function so that the subsequent code is not executed
	if r.Method != "POST" {
		// Use the Header().Set() method to add an 'Allow: POST' header to the
		// response header map. The first parameter is the header name, and
		// the second parameter is the header value. Lets	 the user know
		// which request methods are supported for that particular URL.
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Use the http.ListenAndServe() function to start a new web server. We pas
	// two parameters: the TCP network address to listen on (in this case ":4000
	// and the servemux we just created. If http.ListenAndServe() returns an er
	// we use the log.Fatal() function to log the error message and exit.

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}