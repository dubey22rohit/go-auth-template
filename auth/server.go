package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *Application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ErrorLog:     log.New(app.Logger, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.Logger.PrintInfo("shutting down auth service", map[string]string{
			"signal": s.String(),
		}, app.Config.Debug)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.Logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		}, app.Config.Debug)

		app.wg.Wait()

		shutdownError <- nil
	}()

	app.Logger.PrintInfo("starting auth service", map[string]string{
		"addr": srv.Addr,
		"env":  app.Config.Env,
	}, app.Config.Debug)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.Logger.PrintInfo("auth service stopped", map[string]string{
		"addr": srv.Addr,
	}, app.Config.Debug)

	return nil
}
