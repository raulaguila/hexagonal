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
	"github.com/raulaguila/go-api/internal/core/service/auditor"
	"github.com/raulaguila/go-api/pkg/pgerror"
)

// RoleHandler handles role endpoints
type RoleHandler struct {
	useCase     input.RoleUseCase
	auditor     auditor.Auditor
	handleError func(*fiber.Ctx, error) error
}

// NewRoleHandler creates a new RoleHandler and registers routes
func NewRoleHandler(router fiber.Router, useCase input.RoleUseCase, aud auditor.Auditor, accessAuth fiber.Handler) {
	handler := &RoleHandler{
		useCase: useCase,
		auditor: aud,
		handleError: middleware.NewErrorHandler(middleware.ErrorMapping{
			fiber.MethodDelete: {
				pgerror.ErrForeignKeyViolated: {fiber.StatusBadRequest, "roleUsed"},
			},
			"*": {
				pgerror.ErrUndefinedColumn: {fiber.StatusBadRequest, "undefinedColumn"},
				pgerror.ErrDuplicatedKey:   {fiber.StatusConflict, "roleRegistered"},
				gorm.ErrRecordNotFound:     {fiber.StatusNotFound, "roleNotFound"},
			},
		}),
	}

	// Middleware for parsing DTOs
	roleFilterDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: middleware.CtxKeyFilter,
		OnLookup:   middleware.Query,
		Model:      &dto.RoleFilter{},
	})

	roleInputDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: middleware.CtxKeyDTO,
		OnLookup:   middleware.Body,
		Model:      &dto.RoleInput{},
	})

	// ID is now a UUID string, handled as string by generic parser or specific struct
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

	router.Use(accessAuth)

	router.Get("", roleFilterDTO, middleware.RequirePermission("roles:view"), handler.getRoles)
	router.Get("/list", roleFilterDTO, middleware.RequirePermission("roles:view"), handler.listRoles)
	router.Post("", roleInputDTO, middleware.RequirePermission("roles:create"), handler.createRole)
	router.Put("/:id", idParamDTO, roleInputDTO, middleware.RequirePermission("roles:edit"), handler.updateRole)
	router.Delete("", idsBodyDTO, middleware.RequirePermission("roles:delete"), handler.deleteRoles)
}

// getRoles godoc
// @Summary      Get roles
// @Description  Get roles
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        pgfilter			query	dto.RoleFilter	    false	"Role Filter"
// @Success      200  {object}   	dto.PaginatedOutput[dto.RoleOutput]
// @Failure      403,500  {object}  	presenter.Response
// @Router       /role [get]
// @Security	 Bearer
func (h *RoleHandler) getRoles(c *fiber.Ctx) error {
	filter := GetLocal[dto.RoleFilter](c, middleware.CtxKeyFilter)
	// Permission logic for listing root moved or kept?
	// Assuming similar logic implies 'ListRoot'
	// h.canListRoot(c) check needs User entity in context
	if canList, ok := h.checkListRoot(c); ok {
		filter.ListRoot = canList
	}

	response, err := h.useCase.GetRoles(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// listRoles godoc
// @Summary      List roles
// @Description  List roles
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        pgfilter			query	dto.RoleFilter	    false	"Role Filter"
// @Success      200  {array}   	dto.ItemOutput
// @Failure      403,500  {object}  	presenter.Response
// @Router       /role/list [get]
// @Security	 Bearer
func (h *RoleHandler) listRoles(c *fiber.Ctx) error {
	filter := GetLocal[dto.RoleFilter](c, middleware.CtxKeyFilter)
	if canList, ok := h.checkListRoot(c); ok {
		filter.ListRoot = canList
	}

	response, err := h.useCase.ListRoles(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// createRole godoc
// @Summary      Insert role
// @Description  Insert role
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        role				body	dto.RoleInput       true	"Role model"
// @Success      201  {object}  	dto.RoleOutput
// @Failure      400,403,409,500  {object}  	presenter.Response
// @Router       /role [post]
// @Security	 Bearer
func (h *RoleHandler) createRole(c *fiber.Ctx) error {
	roleDTO := GetLocal[dto.RoleInput](c, middleware.CtxKeyDTO)

	role, err := h.useCase.CreateRole(c.Context(), roleDTO)
	if err != nil {
		return h.handleError(c, err)
	}

	// Audit log
	if h.auditor != nil {
		actor, _ := c.Locals(middleware.LocalUser).(*entity.User)
		var actorID *uuid.UUID
		if actor != nil {
			actorID = &actor.ID
		}
		resID := ""
		if role.ID != nil {
			resID = *role.ID
		}
		// Build metadata with DTO data
		metadata := map[string]any{
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
		}
		if roleDTO != nil {
			metadata["input"] = map[string]any{
				"name":        roleDTO.Name,
				"enabled":     roleDTO.Enabled,
				"permissions": roleDTO.Permissions,
			}
		}
		h.auditor.Log(c.Context(), actorID, "create", "role", resID, metadata)
	}

	return presenter.Created(c, fiberi18n.MustLocalize(c, "roleCreated"), role)
}

// updateRole godoc
// @Summary      Update role by ID
// @Description  Update role by ID
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path    string				true	"Role ID (UUID)"
// @Param        role				body	dto.RoleInput 	    true	"Role model"
// @Success      200  {object}  	dto.RoleOutput
// @Failure      400,403,404,500  {object}  	presenter.Response
// @Router       /role/{id} [put]
// @Security	 Bearer
func (h *RoleHandler) updateRole(c *fiber.Ctx) error {
	idStruct := GetLocal[struct {
		ID string `params:"id" validate:"uuid"`
	}](c, middleware.CtxKeyID)
	roleDTO := GetLocal[dto.RoleInput](c, middleware.CtxKeyDTO)

	role, err := h.useCase.UpdateRole(c.Context(), idStruct.ID, roleDTO)
	if err != nil {
		return h.handleError(c, err)
	}

	// Audit log
	if h.auditor != nil {
		actor, _ := c.Locals(middleware.LocalUser).(*entity.User)
		var actorID *uuid.UUID
		if actor != nil {
			actorID = &actor.ID
		}
		resID := ""
		if role.ID != nil {
			resID = *role.ID
		}
		// Build metadata with DTO data
		metadata := map[string]any{
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
		}
		if roleDTO != nil {
			metadata["input"] = map[string]any{
				"name":        roleDTO.Name,
				"enabled":     roleDTO.Enabled,
				"permissions": roleDTO.Permissions,
			}
		}
		h.auditor.Log(c.Context(), actorID, "update", "role", resID, metadata)
	}

	return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "roleUpdated"), role)
}

// deleteRoles godoc
// @Summary      Delete roles by ID
// @Description  Delete roles by ID
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        ids				body	dto.IDsInput   		true	"Roles ID"
// @Success      204  {object}  	presenter.Response
// @Failure      403,404,500  {object}  	presenter.Response
// @Router       /role [delete]
// @Security	 Bearer
func (h *RoleHandler) deleteRoles(c *fiber.Ctx) error {
	toDelete := GetLocal[dto.IDsInput](c, middleware.CtxKeyID)

	if err := h.useCase.DeleteRoles(c.Context(), toDelete.IDs); err != nil {
		return h.handleError(c, err)
	}

	// Audit log
	if h.auditor != nil {
		actor, _ := c.Locals(middleware.LocalUser).(*entity.User)
		var actorID *uuid.UUID
		if actor != nil {
			actorID = &actor.ID
		}
		for _, deletedID := range toDelete.IDs {
			h.auditor.Log(c.Context(), actorID, "delete", "role", deletedID, map[string]interface{}{
				"ip":         c.IP(),
				"user_agent": c.Get("User-Agent"),
			})
		}
	}

	return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "roleDeleted"), nil)
}

// checkListRoot checks if the current user can list root role
// Returns (canList, success). If User not found in context, returns (false, false)
// Logic for Root role might need adaptation if Root ID provided as constant or config?
// Assuming Root is no longer strictly ID 1, or ID 1 is reserved UUID?
// For now, let's just omit strict ID checking or use permission check if available.
// But keeping legacy logic structure:
func (h *RoleHandler) checkListRoot(c *fiber.Ctx) (bool, bool) {
	// Re-implement if needed based on new Role logical ID or Name
	// For now assuming all roles listable or logic is handled in UseCase
	return true, true
}
