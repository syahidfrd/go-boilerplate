package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appMiddleware "github.com/syahidfrd/go-boilerplate/delivery/middleware"
	"github.com/syahidfrd/go-boilerplate/entity"
	"github.com/syahidfrd/go-boilerplate/mocks"
)

func TestGenerateCID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	mockLogger := new(mocks.Logger)
	cid := appMiddleware.NewMiddleware(mockLogger).RequestID()
	h := cid(handler)
	err := h(c)

	require.NoError(t, err)
	assert.NotNil(t, rec.Header().Get(entity.RequestIDHeader))
}
