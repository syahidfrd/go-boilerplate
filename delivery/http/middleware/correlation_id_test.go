package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syahidfrd/go-boilerplate/delivery/http/middleware"
)

func TestGenerateCorrelationID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	cid := middleware.GenerateCorrelationID()
	h := cid(handler)
	err := h(c)

	require.NoError(t, err)
	assert.NotNil(t, rec.Header().Get(middleware.HeaderXCorrelationID))
}
