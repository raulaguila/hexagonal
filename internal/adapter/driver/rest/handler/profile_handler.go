package handler

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/adapter/driver/rest/middleware"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/pkg/pgerror"
)

const (
	localDTO    = "localDTO"
	localFilter = "localFilter"
	localID     = "localID"
	paramID     = "id"
)

// ProfileHandler handles profile endpoints
type ProfileHandler struct {
	useCase     input.ProfileUseCase
	handleError func(*fiber.Ctx, error) error
}

// NewProfileHandler creates a new ProfileHandler and registers routes
func NewProfileHandler(router fiber.Router, useCase input.ProfileUseCase, accessAuth fiber.Handler) {
	handler := &ProfileHandler{
		useCase: useCase,
		handleError: middleware.NewErrorHandler(middleware.ErrorMapping{
			fiber.MethodDelete: {
				pgerror.ErrForeignKeyViolated: {fiber.StatusBadRequest, "profileUsed"},
			},
			"*": {
				pgerror.ErrUndefinedColumn: {fiber.StatusBadRequest, "undefinedColumn"},
				pgerror.ErrDuplicatedKey:   {fiber.StatusConflict, "profileRegistered"},
				gorm.ErrRecordNotFound:     {fiber.StatusNotFound, "profileNotFound"},
			},
		}),
	}

	// Middleware for parsing DTOs
	profileFilterDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: localFilter,
		OnLookup:   middleware.Query,
		Model:      &dto.ProfileFilter{},
	})

	profileInputDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: localDTO,
		OnLookup:   middleware.Body,
		Model:      &dto.ProfileInput{},
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

	router.Use(accessAuth)

	router.Get("", profileFilterDTO, handler.getProfiles)
	router.Get("/list", profileFilterDTO, handler.listProfiles)
	router.Post("", profileInputDTO, handler.createProfile)
	router.Put("/:"+paramID, idParamDTO, profileInputDTO, handler.updateProfile)
	router.Delete("", idsBodyDTO, handler.deleteProfiles)
}

// getProfiles godoc
// @Summary      Get profiles
// @Description  Get profiles
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        pgfilter			query	dto.ProfileFilter	false	"Profile Filter"
// @Success      200  {object}   	dto.PaginatedOutput[dto.ProfileOutput]
// @Failure      500  {object}  	presenter.Response
// @Router       /profile [get]
// @Security	 Bearer
func (h *ProfileHandler) getProfiles(c *fiber.Ctx) error {
	filter := c.Locals(localFilter).(*dto.ProfileFilter)
	filter.ListRoot = h.canListRoot(c)

	response, err := h.useCase.GetProfiles(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// listProfiles godoc
// @Summary      List profiles
// @Description  List profiles
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        pgfilter			query	dto.ProfileFilter	false	"Profile Filter"
// @Success      200  {array}   	dto.ItemOutput
// @Failure      500  {object}  	presenter.Response
// @Router       /profile/list [get]
// @Security	 Bearer
func (h *ProfileHandler) listProfiles(c *fiber.Ctx) error {
	filter := c.Locals(localFilter).(*dto.ProfileFilter)
	filter.ListRoot = h.canListRoot(c)

	response, err := h.useCase.ListProfiles(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// createProfile godoc
// @Summary      Insert profile
// @Description  Insert profile
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        profile			body	dto.ProfileInput	true	"Profile model"
// @Success      201  {object}  	dto.ProfileOutput
// @Failure      400  {object}  	presenter.Response
// @Failure      409  {object}  	presenter.Response
// @Failure      500  {object}  	presenter.Response
// @Router       /profile [post]
// @Security	 Bearer
func (h *ProfileHandler) createProfile(c *fiber.Ctx) error {
	profileDTO := c.Locals(localDTO).(*dto.ProfileInput)

	profile, err := h.useCase.CreateProfile(c.Context(), profileDTO)
	if err != nil {
		return h.handleError(c, err)
	}

	return presenter.Created(c, fiberi18n.MustLocalize(c, "profileCreated"), profile)
}

// updateProfile godoc
// @Summary      Update profile by ID
// @Description  Update profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path    uint				true	"Profile ID"
// @Param        profile			body	dto.ProfileInput 	true	"Profile model"
// @Success      200  {object}  	dto.ProfileOutput
// @Failure      400  {object}  	presenter.Response
// @Failure      404  {object}  	presenter.Response
// @Failure      500  {object}  	presenter.Response
// @Router       /profile/{id} [put]
// @Security	 Bearer
func (h *ProfileHandler) updateProfile(c *fiber.Ctx) error {
	id := c.Locals(localID).(*struct {
		ID uint `params:"id"`
	})
	profileDTO := c.Locals(localDTO).(*dto.ProfileInput)

	profile, err := h.useCase.UpdateProfile(c.Context(), id.ID, profileDTO)
	if err != nil {
		return h.handleError(c, err)
	}

	return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "profileUpdated"), profile)
}

// deleteProfiles godoc
// @Summary      Delete profiles by ID
// @Description  Delete profiles by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        ids				body	dto.IDsInput   		true	"Profiles ID"
// @Success      204  {object}  	presenter.Response
// @Failure      404  {object}  	presenter.Response
// @Failure      500  {object}  	presenter.Response
// @Router       /profile [delete]
// @Security	 Bearer
func (h *ProfileHandler) deleteProfiles(c *fiber.Ctx) error {
	toDelete := c.Locals(localID).(*dto.IDsInput)

	if err := h.useCase.DeleteProfiles(c.Context(), toDelete.IDs); err != nil {
		return h.handleError(c, err)
	}

	return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "profileDeleted"), nil)
}

// canListRoot checks if the current user can list root profile
func (h *ProfileHandler) canListRoot(c *fiber.Ctx) bool {
	if user, ok := c.Locals(middleware.LocalUser).(*entity.User); ok && user != nil && user.Auth != nil {
		return user.Auth.ProfileID == 1
	}
	return false
}
