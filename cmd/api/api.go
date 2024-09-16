package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	addr string
}

type application struct {
	config config
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello World!"))
	})

	return r
}

func (app *application) run(handler http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINFO, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown

	if err != nil {
		return err
	}

	return nil
}
