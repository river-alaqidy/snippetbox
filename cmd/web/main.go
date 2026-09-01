package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"log/slog" // TODO: look into https://pkg.go.dev/go.uber.org/zap
	"net/http"
	"os"

	"snippetbox.riveralaqidy.net/internal/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// application struct to hold application-wide dependencies
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

// openDB return db connection pool for a given dsn (data source name)
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

func main() {
	// new command-line flag 'addr', default value ":4000", and short description
	// value of flag is stored in addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// define a new command-line flag for the MySQL DSN string
	dsn := flag.String("dsn", os.Getenv("DSN"), "MySQL data source name")
	// read command-line flag and assign it to addr variable
	// call before use addr variable otherwise it will always contian the default value
	flag.Parse()

	// initialize new structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// helper to create connection pool
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// defer closing of db connection before the main function ends
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// init a new instance of application stuct containing dependencies
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
