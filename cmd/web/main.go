package main

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"ssnipp.com/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	debug          bool
	logger         *slog.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	allowSignup    bool
}

// The main() function, which is the entry point for the application.
func main() {
	// Initialize a new logger instance...
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load the .env file into the environment...
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
		os.Exit(1)
	}

	// Read the PORT environment variable to get the address that the server
	// should listen on. If the environment variable isn't set, we default to
	// ":4000", which means the server will listen on all incoming HTTP requests
	// on port 4000.
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":4000"
	}

	// Read the DEBUG environment variable to determine whether the application
	// should run in debug mode. If the environment variable isn't set, we default
	// to false.
	debugStr := os.Getenv("DEBUG")
	if debugStr == "" {
		debugStr = "false"
	}

	// Parse the DEBUG environment variable to a boolean value...
	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		logger.Error("Error parsing DEBUG environment variable")
		os.Exit(1)
	}

	// Read the DB_USERNAME environment variable. If it's not set, log an error
	dbUser := os.Getenv("DB_USERNAME")
	if dbUser == "" {
		logger.Error("DB_USERNAME environment variable not set")
		os.Exit(1)
	}

	// Read the DB_PASSWORD environment variable. If it's not set, log an error
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		logger.Error("DB_PASSWORD environment variable not set")
		os.Exit(1)
	}

	// Read the DB_DATABASE environment variable. If it's not set, log an error
	dbDatabase := os.Getenv("DB_DATABASE")
	if dbDatabase == "" {
		logger.Error("DB_DATABASE environment variable not set")
		os.Exit(1)
	}

	// Construct a DSN from the DB_USERNAME, DB_PASSWORD and DB_DATABASE environment
	dsn := dbUser + ":" + dbPassword + "@/" + dbDatabase + "?parseTime=true"

	// Open a connection to the database...
	db, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Read the ALLOW_SIGNUP environment variable to determine whether new user
	// signups should be allowed. If the environment variable isn't set, we default
	// to true.
	allowSignupStr := os.Getenv("ALLOW_SIGNUP")
	if allowSignupStr == "" {
		allowSignupStr = "true"
	}

	// Parse the ALLOW_SIGNUP environment variable to a boolean value...
	allowSignup, err := strconv.ParseBool(allowSignupStr)
	if err != nil {
		logger.Error("Error parsing ALLOW_SIGNUP environment variable")
		os.Exit(1)
	}

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	// Initialize a new session manager...
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize a new application instance...
	app := &application{
		debug:          debug,
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		allowSignup:    allowSignup,
	}

	// Initialize a new HTTP server...
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", "addr", srv.Addr)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function opens a connection to the MySQL database.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
