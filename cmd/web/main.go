package main

import (
	"flag"
	"log/slog" // TODO: look into https://pkg.go.dev/go.uber.org/zap
	"net/http"
	"os"
)

// application struct to hold application-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {
	// new command-line flag 'addr', default value ":4000", and short description
	// value of flag is stored in addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")
	// read command-line flage and assign it to addr variable
	// call before use addr variable otherwise it will always contian the default value
	flag.Parse()

	// initialize new structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// init a new instance of application stuct containing dependencies
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", "addr", *addr)

	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
