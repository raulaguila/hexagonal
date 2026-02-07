package handler

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/adapter/driver/rest/middleware"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/pkg/pgerror"
)

// UserHandler handles user endpoints
type UserHandler struct {
	useCase     input.UserUseCase
	auditor     input.AuditorUseCase
	handleError func(*fiber.Ctx, error) error
}

// NewUserHandler creates a new UserHandler and registers routes
func NewUserHandler(router fiber.Router, useCase input.UserUseCase, aud input.AuditorUseCase, accessAuth fiber.Handler) {
	handler := &UserHandler{
		useCase: useCase,
		auditor: aud,
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
		ContextKey: middleware.CtxKeyFilter,
		OnLookup:   middleware.Query,
		Model:      &dto.UserFilter{},
	})

	userInputDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: middleware.CtxKeyDTO,
		OnLookup:   middleware.Body,
		Model:      &dto.UserInput{},
	})

	passwordInputDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: middleware.CtxKeyDTO,
		OnLookup:   middleware.Body,
		Model:      &dto.PasswordInput{},
	})

	idParamDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: middleware.CtxKeyID,
		OnLookup:   middleware.Params,
		Model: &struct {
			ID string `params:"id" validate:"uuid"`
		}{},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return presenter.BadRequest(c, "invalidID")
		},
	})

	idsBodyDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: middleware.CtxKeyID,
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
	router.Delete("/pass", middleware.RequirePermission("users:edit"), handler.resetUserPassword)
	router.Get("", userFilterDTO, middleware.RequirePermission("users:view"), handler.getUsers)
	router.Post("", userInputDTO, middleware.RequirePermission("users:create"), handler.createUser)
	router.Put("/:id", idParamDTO, userInputDTO, middleware.RequirePermission("users:edit"), handler.updateUser)
	router.Delete("", idsBodyDTO, middleware.RequirePermission("users:delete"), handler.deleteUser)
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
// @Failure      403,500  {object}  	presenter.Response
// @Router       /user [get]
// @Security	 Bearer
func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	filter := GetLocal[dto.UserFilter](c, middleware.CtxKeyFilter)

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
// @Failure      400,403,409,500  {object}  	presenter.Response
// @Router       /user [post]
// @Security	 Bearer
func (h *UserHandler) createUser(c *fiber.Ctx) error {
	userDTO := GetLocal[dto.UserInput](c, middleware.CtxKeyDTO)

	user, err := h.useCase.CreateUser(c.Context(), userDTO)
	if err != nil {
		return h.handleError(c, err)
	}

	// Audit log
	go func(c *fiber.Ctx, input *dto.UserInput, output *dto.UserOutput) {
		if h.auditor != nil {
			actor, _ := c.Locals(middleware.LocalUser).(*entity.User)
			var actorID *uuid.UUID
			if actor != nil {
				actorID = &actor.ID
			}
			resID := ""
			if output.ID != nil {
				resID = *output.ID
			}
			// Build metadata with DTO data (mask sensitive fields)
			metadata := map[string]any{
				"ip":         c.IP(),
				"user_agent": c.Get("User-Agent"),
			}

			if input != nil {
				metadata["input"] = map[string]any{
					"name":     input.Name,
					"username": input.Username,
					"email":    input.Email,
					"status":   input.Status,
					"role_ids": input.RoleIDs,
				}
			}
			h.auditor.Log(actorID, "create", "user", resID, metadata)
		}
	}(c, userDTO, user)

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
// @Param        id					path		string				true	"User ID"
// @Param        user				body		dto.UserInput		true	"User model"
// @Success      200  {object}  	dto.UserOutput
// @Failure      400,403,404,500  {object}  	presenter.Response
// @Router       /user/{id} [put]
// @Security	 Bearer
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	idStruct := GetLocal[struct {
		ID string `params:"id" validate:"uuid"`
	}](c, middleware.CtxKeyID)
	userDTO := GetLocal[dto.UserInput](c, middleware.CtxKeyDTO)

	user, err := h.useCase.UpdateUser(c.Context(), idStruct.ID, userDTO)
	if err != nil {
		return h.handleError(c, err)
	}

	// Audit log
	go func(c *fiber.Ctx, input *dto.UserInput, output *dto.UserOutput) {
		if h.auditor != nil {
			actor, _ := c.Locals(middleware.LocalUser).(*entity.User)
			var actorID *uuid.UUID
			if actor != nil {
				actorID = &actor.ID
			}
			resID := ""
			if output.ID != nil {
				resID = *output.ID
			}
			// Build metadata with DTO data (mask sensitive fields)
			metadata := map[string]any{
				"ip":         c.IP(),
				"user_agent": c.Get("User-Agent"),
			}
			if input != nil {
				metadata["input"] = map[string]any{
					"name":     input.Name,
					"username": input.Username,
					"email":    input.Email,
					"status":   input.Status,
					"role_ids": input.RoleIDs,
				}
			}
			h.auditor.Log(actorID, "update", "user", resID, metadata)
		}
	}(c, userDTO, user)

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
// @Failure      403,404,500  {object}  	presenter.Response
// @Router       /user [delete]
// @Security	 Bearer
func (h *UserHandler) deleteUser(c *fiber.Ctx) error {
	toDelete := GetLocal[dto.IDsInput](c, middleware.CtxKeyID)

	if err := h.useCase.DeleteUsers(c.Context(), toDelete.IDs); err != nil {
		return h.handleError(c, err)
	}

	// Audit log
	go func(c *fiber.Ctx, input *dto.IDsInput) {
		if h.auditor != nil {
			actor, _ := c.Locals(middleware.LocalUser).(*entity.User)
			var actorID *uuid.UUID
			if actor != nil {
				actorID = &actor.ID
			}
			for _, deletedID := range input.IDs {
				h.auditor.Log(actorID, "delete", "user", deletedID, map[string]any{
					"ip":         c.IP(),
					"user_agent": c.Get("User-Agent"),
				})
			}
		}
	}(c, toDelete)

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
// @Failure      403,404,500  {object}  	presenter.Response
// @Router       /user/pass [delete]
// @Security	 Bearer
func (h *UserHandler) resetUserPassword(c *fiber.Ctx) error {
	// GetQuery helper might not be defined in this file, but assuming it exists as per previous file content
	// But previous file content didn't show GetQuery definition, it showed call.
	// It's likely in `base.go` or similar in same package. I'll invoke it as `GetQuery`.
	// Wait, earlier I saw `GetLocal` usage which relied on external definition or local helper?
	// In `userHandler` line 202: `email, err := GetQuery(c, "email")`.
	// I need to make sure I don't break existing helper usage.

	// I'll assume GetQuery is available in package handler (likely in base.go)

	email, err := GetQuery(c, "email")
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
// @Failure      403,404,500  {object}  	presenter.Response
// @Router       /user/pass [put]
// @Security	 Bearer
func (h *UserHandler) setUserPassword(c *fiber.Ctx) error {
	email, err := GetQuery(c, "email")
	if err != nil {
		return h.handleError(c, err)
	}

	pass := GetLocal[dto.PasswordInput](c, middleware.CtxKeyDTO)

	if err := h.useCase.SetPassword(c.Context(), email, pass); err != nil {
		return h.handleError(c, err)
	}

	return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "passSet"), nil)
}
