package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

// Config holds the module configuration
type Config struct {
	ModuleName      string
	LowerModuleName string
	FirstCharLower  string
	BasePath        string
}

func main() {
	moduleName := flag.String("name", "", "Name of the module (e.g. Product)")
	flag.Parse()

	if *moduleName == "" {
		fmt.Println("Please provide a module name using -name flag")
		os.Exit(1)
	}

	// Normalize name (ensure PascalCase for struct names)
	name := *moduleName
	if len(name) > 0 {
		r := []rune(name)
		r[0] = unicode.ToUpper(r[0])
		name = string(r)
	}

	lowerName := strings.ToLower(name)
	firstLower := string(unicode.ToLower([]rune(name)[0])) + name[1:]

	cfg := Config{
		ModuleName:      name,
		LowerModuleName: lowerName,
		FirstCharLower:  firstLower,
		BasePath:        ".", // Execute from project root
	}

	fmt.Printf("🚀 Scaffolding module '%s'...\n", name)

	files := map[string]string{
		// Domain Entity
		fmt.Sprintf("internal/core/domain/entity/%s.go", lowerName): entityTemplate,

		// DTOs
		fmt.Sprintf("internal/core/dto/%s.go", lowerName): dtoTemplate,

		// Ports
		fmt.Sprintf("internal/core/port/input/%s_usecase.go", lowerName):     inputPortTemplate,
		fmt.Sprintf("internal/core/port/output/%s_repository.go", lowerName): outputPortTemplate,

		// UseCase Implementation
		fmt.Sprintf("internal/core/usecase/%s/%s_usecase.go", lowerName, lowerName): useCaseTemplate,

		// Persistence (Postgres)
		fmt.Sprintf("internal/adapter/driven/persistence/postgres/model/%s.go", lowerName):                 modelTemplate,
		fmt.Sprintf("internal/adapter/driven/persistence/postgres/repository/%s_repository.go", lowerName): repoTemplate,

		// Handler
		fmt.Sprintf("internal/adapter/driver/rest/handler/%sHandler.go", lowerName): handlerTemplate,
	}

	for path, tmpl := range files {
		if err := generateFile(path, tmpl, cfg); err != nil {
			fmt.Printf("❌ Error generating %s: %v\n", path, err)
		} else {
			fmt.Printf("✅ Generated %s\n", path)
		}
	}

	fmt.Println("\n🎉 Scaffolding complete! Don't forget to:")
	fmt.Printf("1. Register the repository in internal/di/container.go\n")
	fmt.Printf("2. Register the usecase in internal/di/container.go\n")
	fmt.Printf("3. Register the handler in internal/adapter/driver/rest/server.go\n")
	fmt.Printf("4. Run 'go mod tidy'\n")
}

func generateFile(path string, tmplContent string, cfg Config) error {
	// Create directories if not exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists")
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New("file").Parse(tmplContent)
	if err != nil {
		return err
	}

	return t.Execute(f, cfg)
}

// --- Templates ---

const entityTemplate = `package entity

import (
	"time"
)

// {{.ModuleName}} represents the domain entity
type {{.ModuleName}} struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	// Add fields here
}

// New{{.ModuleName}} creates a new {{.ModuleName}}
func New{{.ModuleName}}() *{{.ModuleName}} {
	now := time.Now()
	return &{{.ModuleName}}{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (e *{{.ModuleName}}) Validate() error {
	// Add validation logic
	return nil
}
`

const dtoTemplate = `package dto

import "time"

// {{.ModuleName}}Input represents data for creating/updating
type {{.ModuleName}}Input struct {
	// Add fields
}

func (i *{{.ModuleName}}Input) Validate() error {
	return nil
}

// {{.ModuleName}}Output represents output data
type {{.ModuleName}}Output struct {
	ID        uint      ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// {{.ModuleName}}Filter search filters
type {{.ModuleName}}Filter struct {
	PaginatedInput
	// Add filters
}
`

const inputPortTemplate = `package input

import (
	"context"

	"github.com/raulaguila/go-api/internal/core/dto"
)

type {{.ModuleName}}UseCase interface {
	Get{{.ModuleName}}s(ctx context.Context, filter *dto.{{.ModuleName}}Filter) (*dto.PaginatedOutput[dto.{{.ModuleName}}Output], error)
	Get{{.ModuleName}}(ctx context.Context, id uint) (*dto.{{.ModuleName}}Output, error)
	Create{{.ModuleName}}(ctx context.Context, input *dto.{{.ModuleName}}Input) (*dto.{{.ModuleName}}Output, error)
	Update{{.ModuleName}}(ctx context.Context, id uint, input *dto.{{.ModuleName}}Input) (*dto.{{.ModuleName}}Output, error)
	Delete{{.ModuleName}}s(ctx context.Context, ids []uint) error
}
`

const outputPortTemplate = `package output

import (
	"context"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
)

type {{.ModuleName}}Repository interface {
	Count(ctx context.Context, filter *dto.{{.ModuleName}}Filter) (int64, error)
	FindAll(ctx context.Context, filter *dto.{{.ModuleName}}Filter) ([]*entity.{{.ModuleName}}, error)
	FindByID(ctx context.Context, id uint) (*entity.{{.ModuleName}}, error)
	Create(ctx context.Context, data *entity.{{.ModuleName}}) error
	Update(ctx context.Context, data *entity.{{.ModuleName}}) error
	Delete(ctx context.Context, ids []uint) error
}
`

const useCaseTemplate = `package {{.LowerModuleName}}

import (
	"context"
	"errors"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/internal/core/port/output"
)

type {{.FirstCharLower}}UseCase struct {
	repo output.{{.ModuleName}}Repository
}

func New{{.ModuleName}}UseCase(repo output.{{.ModuleName}}Repository) input.{{.ModuleName}}UseCase {
	return &{{.FirstCharLower}}UseCase{repo: repo}
}

func (uc *{{.FirstCharLower}}UseCase) Get{{.ModuleName}}s(ctx context.Context, filter *dto.{{.ModuleName}}Filter) (*dto.PaginatedOutput[dto.{{.ModuleName}}Output], error) {
	count, err := uc.repo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}
	
	items, err := uc.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	outputs := make([]dto.{{.ModuleName}}Output, len(items))
	for i, item := range items {
		outputs[i] = dto.{{.ModuleName}}Output{
			ID: item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return dto.NewPaginatedOutput(outputs, filter.Page, filter.Limit, count), nil
}

func (uc *{{.FirstCharLower}}UseCase) Get{{.ModuleName}}(ctx context.Context, id uint) (*dto.{{.ModuleName}}Output, error) {
	item, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return &dto.{{.ModuleName}}Output{
		ID: item.ID,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (uc *{{.FirstCharLower}}UseCase) Create{{.ModuleName}}(ctx context.Context, input *dto.{{.ModuleName}}Input) (*dto.{{.ModuleName}}Output, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	item := entity.New{{.ModuleName}}()
	// Map input to entity

	if err := uc.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return &dto.{{.ModuleName}}Output{
		ID: item.ID,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (uc *{{.FirstCharLower}}UseCase) Update{{.ModuleName}}(ctx context.Context, id uint, input *dto.{{.ModuleName}}Input) (*dto.{{.ModuleName}}Output, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	item, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	// item.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return &dto.{{.ModuleName}}Output{
		ID: item.ID,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (uc *{{.FirstCharLower}}UseCase) Delete{{.ModuleName}}s(ctx context.Context, ids []uint) error {
	return uc.repo.Delete(ctx, ids)
}
`

const modelTemplate = `package model

import (
	"time"

	"gorm.io/gorm"
)

type {{.ModuleName}}Model struct {
	ID        uint           ` + "`gorm:\"primarykey\"`" + `
	CreatedAt time.Time      ` + "`gorm:\"autoCreateTime\"`" + `
	UpdatedAt time.Time      ` + "`gorm:\"autoUpdateTime\"`" + `
	DeletedAt gorm.DeletedAt ` + "`gorm:\"index\"`" + `
}

func ({{.ModuleName}}Model) TableName() string {
	return "s_{{.LowerModuleName}}"
}
`

const repoTemplate = `package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/output"
)

type {{.LowerModuleName}}Repository struct {
	db *gorm.DB
}

func New{{.ModuleName}}Repository(db *gorm.DB) output.{{.ModuleName}}Repository {
	return &{{.LowerModuleName}}Repository{db: db}
}

func (r *{{.LowerModuleName}}Repository) Count(ctx context.Context, filter *dto.{{.ModuleName}}Filter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.{{.ModuleName}}Model{})
	// Apply filters
	err := query.Count(&count).Error
	return count, err
}

func (r *{{.LowerModuleName}}Repository) FindAll(ctx context.Context, filter *dto.{{.ModuleName}}Filter) ([]*entity.{{.ModuleName}}, error) {
	var models []model.{{.ModuleName}}Model
	query := r.db.WithContext(ctx).Model(&model.{{.ModuleName}}Model{})
	// Apply filters and pagination
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]*entity.{{.ModuleName}}, len(models))
	for i, m := range models {
		entities[i] = &entity.{{.ModuleName}}{
			ID: m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	return entities, nil
}

func (r *{{.LowerModuleName}}Repository) FindByID(ctx context.Context, id uint) (*entity.{{.ModuleName}}, error) {
	var m model.{{.ModuleName}}Model
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &entity.{{.ModuleName}}{
		ID: m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *{{.LowerModuleName}}Repository) Create(ctx context.Context, data *entity.{{.ModuleName}}) error {
	m := &model.{{.ModuleName}}Model{
		// Map fields
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	data.ID = m.ID
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *{{.LowerModuleName}}Repository) Update(ctx context.Context, data *entity.{{.ModuleName}}) error {
	m := &model.{{.ModuleName}}Model{
		ID: data.ID,
		// Map fields
	}
	return r.db.WithContext(ctx).Updates(m).Error
}

func (r *{{.LowerModuleName}}Repository) Delete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&model.{{.ModuleName}}Model{}, ids).Error
}
`

const handlerTemplate = `package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest/middleware"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
)

type {{.ModuleName}}Handler struct {
	useCase     input.{{.ModuleName}}UseCase
	handleError func(*fiber.Ctx, error) error
}

func New{{.ModuleName}}Handler(router fiber.Router, useCase input.{{.ModuleName}}UseCase, auth fiber.Handler) {
	handler := &{{.ModuleName}}Handler{
		useCase: useCase,
		handleError: middleware.NewErrorHandler(middleware.ErrorMapping{
			// Map specific errors if needed
		}),
	}
	
	// Middlewares
	filterDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: middleware.CtxKeyFilter,
		OnLookup:   middleware.Query,
		Model:      &dto.{{.ModuleName}}Filter{},
	})
	
	inputDTO := middleware.ParseDTO(middleware.DTOConfig{
		ContextKey: middleware.CtxKeyDTO,
		OnLookup:   middleware.Body,
		Model:      &dto.{{.ModuleName}}Input{},
	})

	router.Use(auth)
	
	router.Get("", filterDTO, handler.list)
	router.Get("/:id", handler.get)
	router.Post("", inputDTO, handler.create)
	router.Put("/:id", inputDTO, handler.update)
	router.Delete("", handler.delete)
}

func (h *{{.ModuleName}}Handler) list(c *fiber.Ctx) error {
	filter := GetLocal[dto.{{.ModuleName}}Filter](c, middleware.CtxKeyFilter)
	res, err := h.useCase.Get{{.ModuleName}}s(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}
	return c.JSON(res)
}

func (h *{{.ModuleName}}Handler) get(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return presenter.BadRequest(c, "invalid id")
	}
	res, err := h.useCase.Get{{.ModuleName}}(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}
	return c.JSON(res)
}

func (h *{{.ModuleName}}Handler) create(c *fiber.Ctx) error {
	input := GetLocal[dto.{{.ModuleName}}Input](c, middleware.CtxKeyDTO)
	res, err := h.useCase.Create{{.ModuleName}}(c.Context(), input)
	if err != nil {
		return h.handleError(c, err)
	}
	return presenter.Created(c, "created", res)
}

func (h *{{.ModuleName}}Handler) update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return presenter.BadRequest(c, "invalid id")
	}
	input := GetLocal[dto.{{.ModuleName}}Input](c, middleware.CtxKeyDTO)
	res, err := h.useCase.Update{{.ModuleName}}(c.Context(), uint(id), input)
	if err != nil {
		return h.handleError(c, err)
	}
	return c.JSON(res)
}

func (h *{{.ModuleName}}Handler) delete(c *fiber.Ctx) error {
	// Parse IDs from body or query
	return c.SendStatus(fiber.StatusNoContent)
}
`
