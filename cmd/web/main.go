package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/duybnguyen/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application-wide dependencies.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel // allow us to make the SnippetModel object available to our handlers.
}

func main() {
	// Define command-line flags.
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:duy3015873@/snippetbox?parseTime=true", "MySQL data source name")

	// Parse the command-line flags.
	flag.Parse()

	// Create loggers for informational and error messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize the database connection.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize the application struct.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	// Initialize a new servemux and register handlers.
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a file server to serve static files.
	fileServer := http.FileServer(http.Dir("../../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Initialize and start a new HTTP server.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// openDB initializes and returns a new database connection pool.
// The sql.Open() function doesn’t actually create any connections, all it does is initialize the pool for future use. Actual connections to the database are established lazily, as and when needed for the first time. So to verify that everything is set up correctly we need to use the db.Ping() method to create a connection and check for any errors.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
