package handler

import (
	"net/url"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/adapter/driver/rest/middleware"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/pkg/pgerror"
)

// UserHandler handles user endpoints
type UserHandler struct {
	useCase     input.UserUseCase
	handleError func(*fiber.Ctx, error) error
}

// NewUserHandler creates a new UserHandler and registers routes
func NewUserHandler(router fiber.Router, useCase input.UserUseCase, accessAuth fiber.Handler) {
	handler := &UserHandler{
		useCase: useCase,
		handleError: middleware.NewErrorHandler(middleware.ErrorMapping{
			fiber.MethodDelete: {
				pgerror.ErrForeignKeyViolated: {fiber.StatusBadRequest, "userUsed"},
			},
			"*": {
				pgerror.ErrUndefinedColumn:    {fiber.StatusBadRequest, "undefinedColumn"},
				pgerror.ErrDuplicatedKey:      {fiber.StatusConflict, "userRegistered"},
				pgerror.ErrForeignKeyViolated: {fiber.StatusNotFound, "itemNotFound"},
				gorm.ErrRecordNotFound:        {fiber.StatusNotFound, "userNotFound"},
			},
		}),
	}

	// Middleware for parsing DTOs
	userFilterDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: localFilter,
		OnLookup:   middleware.Query,
		Model:      &dto.UserFilter{},
	})

	userInputDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: localDTO,
		OnLookup:   middleware.Body,
		Model:      &dto.UserInput{},
	})

	passwordInputDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: localDTO,
		OnLookup:   middleware.Body,
		Model:      &dto.PasswordInput{},
	})

	idParamDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: localID,
		OnLookup:   middleware.Params,
		Model: &struct {
			ID uint `params:"id"`
		}{},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return presenter.BadRequest(c, "invalidID")
		},
	})

	idsBodyDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: localID,
		OnLookup:   middleware.Body,
		Model:      &dto.IDsInput{},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return presenter.BadRequest(c, "invalidID")
		},
	})

	// Public route for setting password
	router.Put("/pass", passwordInputDTO, handler.setUserPassword)

	// Protected routes
	router.Use(accessAuth)
	router.Delete("/pass", handler.resetUserPassword)
	router.Get("", userFilterDTO, handler.getUsers)
	router.Post("", userInputDTO, handler.createUser)
	router.Put("/:"+paramID, idParamDTO, userInputDTO, handler.updateUser)
	router.Delete("", idsBodyDTO, handler.deleteUser)
}

// getUsers godoc
// @Summary      Get users
// @Description  Get all users
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        pgfilter			query		dto.UserFilter		false	"Optional Filter"
// @Success      200  {object}   	dto.PaginatedOutput[dto.UserOutput]
// @Failure      500  {object}  	presenter.Response
// @Router       /user [get]
// @Security	 Bearer
func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	filter := c.Locals(localFilter).(*dto.UserFilter)

	response, err := h.useCase.GetUsers(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// createUser godoc
// @Summary      Insert user
// @Description  Insert user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        user				body		dto.UserInput		true	"User model"
// @Success      201  {object}  	dto.UserOutput
// @Failure      400  {object}  	presenter.Response
// @Failure      409  {object}  	presenter.Response
// @Failure      500  {object} 		presenter.Response
// @Router       /user [post]
// @Security	 Bearer
func (h *UserHandler) createUser(c *fiber.Ctx) error {
	userDTO := c.Locals(localDTO).(*dto.UserInput)

	user, err := h.useCase.CreateUser(c.Context(), userDTO)
	if err != nil {
		return h.handleError(c, err)
	}

	return presenter.Created(c, fiberi18n.MustLocalize(c, "userCreated"), user)
}

// updateUser godoc
// @Summary      Update user by ID
// @Description  Update user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path		uint				true	"User ID"
// @Param        user				body		dto.UserInput		true	"User model"
// @Success      200  {object}  	dto.UserOutput
// @Failure      400  {object}  	presenter.Response
// @Failure      404  {object}  	presenter.Response
// @Failure      500  {object}  	presenter.Response
// @Router       /user/{id} [put]
// @Security	 Bearer
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	id := c.Locals(localID).(*struct {
		ID uint `params:"id"`
	})
	userDTO := c.Locals(localDTO).(*dto.UserInput)

	user, err := h.useCase.UpdateUser(c.Context(), id.ID, userDTO)
	if err != nil {
		return h.handleError(c, err)
	}

	return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "userUpdated"), user)
}

// deleteUser godoc
// @Summary      Delete user by ID
// @Description  Delete user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool					false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string					false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					body		dto.IDsInput			true	"User ID"
// @Success      204  {object}  	nil
// @Failure      404  {object}  	presenter.Response
// @Failure      500  {object}  	presenter.Response
// @Router       /user [delete]
// @Security	 Bearer
func (h *UserHandler) deleteUser(c *fiber.Ctx) error {
	toDelete := c.Locals(localID).(*dto.IDsInput)

	if err := h.useCase.DeleteUsers(c.Context(), toDelete.IDs); err != nil {
		return h.handleError(c, err)
	}

	return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "userDeleted"), nil)
}

// resetUserPassword godoc
// @Summary      Reset user password by ID
// @Description  Reset user password by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        email				query		string				true 	"User email"
// @Success      200  {object}  	nil
// @Failure      404  {object}  	presenter.Response
// @Failure      500  {object}  	presenter.Response
// @Router       /user/pass [delete]
// @Security	 Bearer
func (h *UserHandler) resetUserPassword(c *fiber.Ctx) error {
	email, err := url.QueryUnescape(c.Query("email", ""))
	if err != nil {
		return h.handleError(c, err)
	}

	if err := h.useCase.ResetPassword(c.Context(), email); err != nil {
		return h.handleError(c, err)
	}

	return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "passReset"), nil)
}

// setUserPassword godoc
// @Summary      Set user password by ID
// @Description  Set user password by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool					false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string					false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        email				query		string					true	"User email" Format(email)
// @Param        password			body		dto.PasswordInput		true	"Password model"
// @Success      200  {object}  	nil
// @Failure      404  {object}  	presenter.Response
// @Failure      500  {object}  	presenter.Response
// @Router       /user/pass [put]
func (h *UserHandler) setUserPassword(c *fiber.Ctx) error {
	email, err := url.QueryUnescape(c.Query("email", ""))
	if err != nil {
		return h.handleError(c, err)
	}

	pass := c.Locals(localDTO).(*dto.PasswordInput)
	if err := pass.Validate(); err != nil {
		return h.handleError(c, err)
	}

	if err := h.useCase.SetPassword(c.Context(), email, pass); err != nil {
		return h.handleError(c, err)
	}

	return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "passSet"), nil)
}
