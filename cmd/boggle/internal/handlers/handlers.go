package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Run starts the webserver.
func Run() error {
	routes()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	server := http.Server{
		Addr: "0.0.0.0:8080",
	}
	go func() {
		log.Info("Listening on http://localhost:8080 ... (Press ctrl + c to quit)")
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "errors starting server")
	case sig := <-shutdown:
		log.Info(fmt.Sprintf("begin shutting down %v", sig))
		shutdownTimeout, _ := time.ParseDuration("5s")
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Info(fmt.Sprintf("Graceful shutdown did not complete in %v : %v", shutdownTimeout, err))
			err = server.Close()
		}

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
