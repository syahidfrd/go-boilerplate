package http

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"github.com/syahidfrd/go-boilerplate/delivery/middleware"
	"github.com/syahidfrd/go-boilerplate/domain"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/utils"
)

type AuthHandler struct {
	AuthUC domain.AuthUsecase
}

// NewAuthHandler will initialize the auth resources endpoint
func NewAuthHandler(e *echo.Echo, middleware *middleware.Middleware, authUC domain.AuthUsecase) {
	handler := &AuthHandler{
		AuthUC: authUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/auth/signup", handler.SignUp)
	apiV1.POST("/auth/signin", handler.SignIn)
}

// SignUp godoc
// @Summary SignUp
// @Description SignUp
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body request.SignUpReq true "SignUp user"
// @Success 200
// @Router /api/v1/auth/signup [post]
func (h *AuthHandler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.SignUpReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.AuthUC.SignUp(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "signup successfully",
	})
}

// SignIn godoc
// @Summary SignIn
// @Description SignIn
// @Tags Auth
// @Accept json
// @Produce json
// @Param signin body request.SignInReq true "SignIn user"
// @Success 200
// @Router /api/v1/auth/signin [post]
func (h *AuthHandler) SignIn(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.SignInReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	accessToken, err := h.AuthUC.SignIn(ctx, &req)

	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"access_token": accessToken,
		},
	})
}
