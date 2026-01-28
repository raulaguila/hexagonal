package handler

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/adapter/driver/rest/middleware"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	useCase     input.AuthUseCase
	handleError func(*fiber.Ctx, error) error
	accessAuth  fiber.Handler
	refreshAuth fiber.Handler
}

// NewAuthHandler creates a new AuthHandler and registers routes
func NewAuthHandler(router fiber.Router, useCase input.AuthUseCase, accessAuth, refreshAuth fiber.Handler) {
	handler := &AuthHandler{
		useCase:     useCase,
		accessAuth:  accessAuth,
		refreshAuth: refreshAuth,
		handleError: middleware.NewErrorHandler(middleware.ErrorMapping{
			"*": {
				gorm.ErrRecordNotFound: {fiber.StatusNotFound, "userNotFound"},
			},
		}),
	}

	router.Post("", handler.login)
	router.Get("", accessAuth, handler.me)
	router.Put("", refreshAuth, handler.refresh)
}

// login godoc
// @Summary      User authentication
// @Description  User authentication
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        credentials		body	dto.LoginInput	true	"Credentials model"
// @Success      200  {object}  	dto.AuthOutput
// @Failure      401  {object}  	presenter.Response
// @Failure      500  {object}  	presenter.Response
// @Router       /auth [post]
func (h *AuthHandler) login(c *fiber.Ctx) error {
	credentials := new(dto.LoginInput)
	if err := c.BodyParser(credentials); err != nil {
		return presenter.BadRequest(c, fiberi18n.MustLocalize(c, "invalidData"))
	}

	authResponse, err := h.useCase.Login(c.Context(), credentials)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(authResponse)
}

// me godoc
// @Summary      User authenticated
// @Description  User authenticated
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Authorization		header	string				false	"User token"
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Success      200  {object}  	dto.UserOutput
// @Failure      401  {object}  	presenter.Response
// @Failure      500  {object}  	presenter.Response
// @Router       /auth [get]
// @Security	 Bearer
func (h *AuthHandler) me(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return presenter.Unauthorized(c, fiberi18n.MustLocalize(c, "unauthorized"))
	}

	user, err := h.useCase.Me(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// refresh godoc
// @Summary      User refresh
// @Description  User refresh
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Authorization		header	string				false	"User token"
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        expire				query	bool				false	"Expire token"
// @Success      200  {object}  	dto.AuthOutput
// @Failure      401  {object}  	presenter.Response
// @Failure      500  {object}  	presenter.Response
// @Router       /auth [put]
func (h *AuthHandler) refresh(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return presenter.Unauthorized(c, fiberi18n.MustLocalize(c, "unauthorized"))
	}

	expire := c.Query("expire", "true") == "true"
	authResponse, err := h.useCase.Refresh(c.Context(), userID, expire)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(authResponse)
}
