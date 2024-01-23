package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	ts, err := NewServer(Config{
		Interface: "127.0.0.1",
		Port:      "",
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	require.NoError(t, err)

	defer func() {
		_ = ts.Stop()
	}()

	errorCh := make(chan error, 2)
	doneCh := make(chan struct{}, 1)

	go func() {
		if err := ts.Run(); err != nil {
			errorCh <- err

			return
		}

		doneCh <- struct{}{}
	}()

	t.Run("connect", func(t *testing.T) {
		require.NotEmpty(t, ts.Address())

		for i := 0; i < 10; i++ {
			rsp, err := http.Get("http://" + ts.Address())
			if err != nil {
				time.Sleep(100 * time.Millisecond)

				continue
			}

			require.Equal(t, http.StatusAccepted, rsp.StatusCode)

			return
		}

		t.Fatal("failed to connect to server")
	})

	t.Run("unbind-ok", func(t *testing.T) {
		require.NoError(t, ts.Unbind())
	})

	t.Run("unbind-err", func(t *testing.T) {
		err := ts.Unbind()

		require.Error(t, err)
		require.Contains(t, err.Error(), "use of closed network connection")
	})

	t.Run("stop-ok", func(t *testing.T) {
		require.NoError(t, ts.Stop())
	})

	t.Run("stop-again-ok", func(t *testing.T) {
		require.NoError(t, ts.Stop())
	})

	select {
	case <-doneCh:
	case err := <-errorCh:
		t.Fatalf("error running server: %s", err)
	case timeout := <-time.After(5 * time.Second):
		t.Fatalf("timeout waiting for server to stop (timeout: %s)", timeout)
	}
}
