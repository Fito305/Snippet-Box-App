package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	// Import the models package. You need to prefix this with
	// Whatever module path you set up. Example:
	// "{your-module_path}/internal/models". You can find it in the go.mod file.
	"snippetbox.felipeacosta.net/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.

// Add a snippet field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
// Add a templateCache field to the application struct.
// Add a formDecoder field to hold a pointer to a form.Decoder instance.
// Add a new sessionManager field to the application struct
// Add a new users field to the application struct.
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Define a new comand-line flag for the MySQL DSN string. When git pushing it change password to web:pass@/
	dsn := flag.String("dsn", "web:Allapattah123456!@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the seperate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close, so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new tamplet cache...
	// And add it to teh application dependencies below
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialize a new session manager. Then we configure it to use our MySQL database
	// as the session store, and set a lifetime of 12 hours (so that sessions atuo exipre 12 hours after first being created.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Make sure that the Secure attribute is set on our session cookies. Setting this means that the cookie will only be sent by a user's web browser when a HTTPS connection is being used (and won't be sent over an unsecure HTTP connection).
	sessionManager.Cookie.Secure = true

	// Initialize a new instance of our application struct, containing the
	// dependencies.
	// Initialize a models.SnippetModel instance and add it to the application dependencies.
	// And add it to the application dependencies.
	// And add the session manager to our application dependencies.
	// Initialize a models.UserModel instance and add it to the application dependencies.
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we want the server to use. In this case the only thing that we're changing is the curve preferences value, so that only elliptic curves with assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Set the server's TLSCOnfig field to use the tlsConfig variable we just created.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting serve on %s", *addr)
	// Because the err variable is now already declared in the code above, we need
	// to use the assignment operator = here, instead of the := `declare and assign` operator.
	// Use the ListenAndServeTLS() method to start the HTTPS server. We pass in the paths to the TLS certificate and corresponding private key as the two parameters.
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool for a given DSN.
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
