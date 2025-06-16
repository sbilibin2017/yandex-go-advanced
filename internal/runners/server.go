package runners

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func RunServer(
	ctx context.Context,
	srv Server,
) error {
	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)

	case err := <-errChan:
		return err
	}
}
