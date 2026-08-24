package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// waitForCallback starts an HTTP server on addr, waits for the OAuth callback,
// validates state, and returns the authorization code.
func waitForCallback(ctx context.Context, addr, expectedState string) (string, error) {
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}

	var shutdownOnce sync.Once
	shutdown := func() {
		shutdownOnce.Do(func() {
			ctx2, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = srv.Shutdown(ctx2)
		})
	}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if errParam := q.Get("error"); errParam != "" {
			desc := q.Get("error_description")
			errCh <- fmt.Errorf("auth error: %s — %s", errParam, desc)
			http.Error(w, "Authentication failed. You may close this tab.", http.StatusBadRequest)
			go shutdown()
			return
		}

		state := q.Get("state")
		if state != expectedState {
			errCh <- fmt.Errorf("state mismatch — possible CSRF attack")
			http.Error(w, "Invalid state.", http.StatusBadRequest)
			go shutdown()
			return
		}

		code := q.Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in callback")
			http.Error(w, "No authorization code received.", http.StatusBadRequest)
			go shutdown()
			return
		}

		// Respond immediately — token exchange happens in the caller.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html><html><body>
<h2>✓ Authenticated</h2>
<p>You can close this tab and return to the terminal.</p>
</body></html>`)

		codeCh <- code
		go shutdown()
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("callback server: %w", err)
		}
	}()

	select {
	case code := <-codeCh:
		return code, nil
	case err := <-errCh:
		return "", err
	case <-ctx.Done():
		shutdown()
		return "", ctx.Err()
	}
}

// signalCallbackDone signals the callback handler that token exchange completed.
func signalCallbackDone(doneCh chan struct{}) {
	select {
	case <-doneCh:
	default:
		close(doneCh)
	}
}
