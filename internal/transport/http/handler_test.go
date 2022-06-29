package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler()

	require.IsType(t, Handler{}, h)
}

func TestNewHandler_Get(t *testing.T) {
	h := NewHandler()

	router := h.GetHTTPHandler()

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/healthcheck")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
