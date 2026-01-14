# 🚀 Guia: Como Adicionar Novas Funcionalidades

> **Passo a passo para adicionar um novo módulo seguindo a Arquitetura Hexagonal**

Este guia usa como exemplo a criação de um módulo de **Produtos**.

---

## 📋 Visão Geral dos Passos

```
1. Entidade     → O que é um Produto?
2. Erros        → O que pode dar errado?
3. DTO          → Quais dados entram e saem?
4. Port Input   → O que a aplicação sabe fazer?
5. Port Output  → O que a aplicação precisa?
6. Repository   → Como salvar no banco?
7. Use Case     → Qual a lógica de negócio?
8. Container    → Como conectar as dependências?
9. Application  → Como expor para interfaces?
10. Handler     → Como receber requisições HTTP?
11. Rotas       → Onde ficam os endpoints?
```

---

## 📦 Passo 1: Criar a Entidade

**Arquivo:** `internal/core/domain/entity/product.go`

```go
package entity

import (
    "time"
)

// Product representa um produto no sistema
type Product struct {
    ID          uint
    Name        string
    Description string
    Price       float64
    Stock       int
    Active      bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// NewProduct cria um novo produto
func NewProduct(name, description string, price float64, stock int) *Product {
    return &Product{
        Name:        name,
        Description: description,
        Price:       price,
        Stock:       stock,
        Active:      true,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
}

// Validate verifica se o produto é válido
func (p *Product) Validate() error {
    if len(p.Name) < 3 {
        return errors.ErrProductNameTooShort
    }
    if p.Price < 0 {
        return errors.ErrProductPriceNegative
    }
    if p.Stock < 0 {
        return errors.ErrProductStockNegative
    }
    return nil
}

// UpdateStock atualiza o estoque
func (p *Product) UpdateStock(quantity int) {
    p.Stock = quantity
    p.UpdatedAt = time.Now()
}

// Deactivate desativa o produto
func (p *Product) Deactivate() {
    p.Active = false
    p.UpdatedAt = time.Now()
}
```

> 💡 **Importante:** A entidade não tem tags GORM! É um objeto puro do domínio.

---

## ❌ Passo 2: Criar os Erros de Domínio

**Arquivo:** `internal/core/domain/errors/product_errors.go`

```go
package errors

import "errors"

// Erros relacionados a produtos
var (
    ErrProductNotFound      = errors.New("product not found")
    ErrProductNameTooShort  = errors.New("product name must have at least 3 characters")
    ErrProductPriceNegative = errors.New("product price cannot be negative")
    ErrProductStockNegative = errors.New("product stock cannot be negative")
    ErrProductInactive      = errors.New("product is inactive")
)
```

---

## 📤 Passo 3: Criar os DTOs

**Arquivo:** `internal/core/dto/product.go`

```go
package dto

import "github.com/raulaguila/go-api/pkg/apperror"

// ProductInput representa os dados para criar/atualizar produto
type ProductInput struct {
    Name        *string  `json:"name"`
    Description *string  `json:"description"`
    Price       *float64 `json:"price"`
    Stock       *int     `json:"stock"`
    Active      *bool    `json:"active"`
}

// Validate valida os dados de entrada
func (p *ProductInput) Validate() error {
    if p.Name != nil && len(*p.Name) < 3 {
        return apperror.InvalidInput("name", "nome deve ter pelo menos 3 caracteres")
    }
    if p.Price != nil && *p.Price < 0 {
        return apperror.InvalidInput("price", "preço não pode ser negativo")
    }
    if p.Stock != nil && *p.Stock < 0 {
        return apperror.InvalidInput("stock", "estoque não pode ser negativo")
    }
    return nil
}

// ProductOutput representa os dados de saída do produto
type ProductOutput struct {
    ID          *uint    `json:"id"`
    Name        *string  `json:"name"`
    Description *string  `json:"description"`
    Price       *float64 `json:"price"`
    Stock       *int     `json:"stock"`
    Active      *bool    `json:"active"`
}

// ProductFilter representa os filtros de busca
type ProductFilter struct {
    PaginationFilter
    Name   string `query:"name"`
    Active *bool  `query:"active"`
}
```

---

## 🔌 Passo 4: Criar o Port de Entrada (Use Case Interface)

**Arquivo:** `internal/core/port/input/product_usecase.go`

```go
package input

import (
    "context"
    
    "github.com/raulaguila/go-api/internal/core/dto"
)

// ProductUseCase define as operações de negócio para produtos
type ProductUseCase interface {
    // GetProducts retorna lista paginada de produtos
    GetProducts(ctx context.Context, filter *dto.ProductFilter) (*dto.PaginatedOutput[dto.ProductOutput], error)
    
    // GetProductByID retorna um produto pelo ID
    GetProductByID(ctx context.Context, id uint) (*dto.ProductOutput, error)
    
    // CreateProduct cria um novo produto
    CreateProduct(ctx context.Context, input *dto.ProductInput) (*dto.ProductOutput, error)
    
    // UpdateProduct atualiza um produto existente
    UpdateProduct(ctx context.Context, id uint, input *dto.ProductInput) (*dto.ProductOutput, error)
    
    // DeleteProducts remove produtos pelos IDs
    DeleteProducts(ctx context.Context, ids []uint) error
}
```

---

## 🔌 Passo 5: Criar o Port de Saída (Repository Interface)

**Arquivo:** `internal/core/port/output/product_repository.go`

```go
package output

import (
    "context"
    
    "github.com/raulaguila/go-api/internal/core/domain/entity"
    "github.com/raulaguila/go-api/internal/core/dto"
)

// ProductRepository define as operações de persistência para produtos
type ProductRepository interface {
    // FindAll retorna todos os produtos com filtros
    FindAll(ctx context.Context, filter *dto.ProductFilter) ([]*entity.Product, error)
    
    // FindByID retorna um produto pelo ID
    FindByID(ctx context.Context, id uint) (*entity.Product, error)
    
    // Count retorna a contagem de produtos
    Count(ctx context.Context, filter *dto.ProductFilter) (int64, error)
    
    // Create cria um novo produto
    Create(ctx context.Context, product *entity.Product) error
    
    // Update atualiza um produto
    Update(ctx context.Context, product *entity.Product) error
    
    // Delete remove produtos pelos IDs
    Delete(ctx context.Context, ids []uint) error
}
```

---

## 💾 Passo 6: Criar o Modelo e Repository (Adapter)

### 6.1 Modelo GORM

**Arquivo:** `internal/adapter/driven/persistence/postgres/model/product_model.go`

```go
package model

import "time"

// ProductModel é o modelo GORM para produtos
type ProductModel struct {
    ID          uint      `gorm:"primarykey"`
    Name        string    `gorm:"column:name;not null"`
    Description string    `gorm:"column:description"`
    Price       float64   `gorm:"column:price;not null"`
    Stock       int       `gorm:"column:stock;default:0"`
    Active      bool      `gorm:"column:active;default:true"`
    CreatedAt   time.Time `gorm:"column:created_at"`
    UpdatedAt   time.Time `gorm:"column:updated_at"`
}

// TableName define o nome da tabela
func (ProductModel) TableName() string {
    return "products"
}
```

### 6.2 Mapper

**Arquivo:** Adicionar em `internal/adapter/driven/persistence/postgres/mapper/mapper.go`

```go
// ProductToModel converte entidade em modelo
func ProductToModel(e *entity.Product) *model.ProductModel {
    if e == nil {
        return nil
    }
    return &model.ProductModel{
        ID:          e.ID,
        Name:        e.Name,
        Description: e.Description,
        Price:       e.Price,
        Stock:       e.Stock,
        Active:      e.Active,
        CreatedAt:   e.CreatedAt,
        UpdatedAt:   e.UpdatedAt,
    }
}

// ProductToEntity converte modelo em entidade
func ProductToEntity(m *model.ProductModel) *entity.Product {
    if m == nil {
        return nil
    }
    return &entity.Product{
        ID:          m.ID,
        Name:        m.Name,
        Description: m.Description,
        Price:       m.Price,
        Stock:       m.Stock,
        Active:      m.Active,
        CreatedAt:   m.CreatedAt,
        UpdatedAt:   m.UpdatedAt,
    }
}

// ProductsToEntities converte slice de modelos para entidades
func ProductsToEntities(models []*model.ProductModel) []*entity.Product {
    return MapSlice(models, ProductToEntity)
}
```

### 6.3 Repository Implementation

**Arquivo:** `internal/adapter/driven/persistence/postgres/repository/product_repository.go`

```go
package repository

import (
    "context"
    
    "gorm.io/gorm"
    
    "github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/mapper"
    "github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
    "github.com/raulaguila/go-api/internal/core/domain/entity"
    "github.com/raulaguila/go-api/internal/core/dto"
    "github.com/raulaguila/go-api/internal/core/port/output"
)

type productRepository struct {
    db *gorm.DB
}

// NewProductRepository cria um novo repository de produtos
func NewProductRepository(db *gorm.DB) output.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) FindAll(ctx context.Context, filter *dto.ProductFilter) ([]*entity.Product, error) {
    var models []*model.ProductModel
    
    query := r.db.WithContext(ctx)
    
    if filter.Name != "" {
        query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
    }
    if filter.Active != nil {
        query = query.Where("active = ?", *filter.Active)
    }
    
    query = query.Offset(filter.GetOffset()).Limit(filter.GetLimit())
    
    if err := query.Find(&models).Error; err != nil {
        return nil, err
    }
    
    return mapper.ProductsToEntities(models), nil
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
    var m model.ProductModel
    if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
        return nil, err
    }
    return mapper.ProductToEntity(&m), nil
}

func (r *productRepository) Count(ctx context.Context, filter *dto.ProductFilter) (int64, error) {
    var count int64
    query := r.db.WithContext(ctx).Model(&model.ProductModel{})
    
    if filter != nil && filter.Name != "" {
        query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
    }
    
    return count, query.Count(&count).Error
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
    m := mapper.ProductToModel(product)
    if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
        return err
    }
    product.ID = m.ID
    return nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
    m := mapper.ProductToModel(product)
    return r.db.WithContext(ctx).Save(m).Error
}

func (r *productRepository) Delete(ctx context.Context, ids []uint) error {
    return r.db.WithContext(ctx).Delete(&model.ProductModel{}, ids).Error
}
```

---

## ⚙️ Passo 7: Criar o Use Case

**Arquivo:** `internal/core/usecase/product/product_usecase.go`

```go
package product

import (
    "context"
    
    "github.com/raulaguila/go-api/internal/core/domain/entity"
    "github.com/raulaguila/go-api/internal/core/domain/errors"
    "github.com/raulaguila/go-api/internal/core/dto"
    "github.com/raulaguila/go-api/internal/core/port/input"
    "github.com/raulaguila/go-api/internal/core/port/output"
)

type productUseCase struct {
    productRepo output.ProductRepository
}

// NewProductUseCase cria um novo use case de produtos
func NewProductUseCase(productRepo output.ProductRepository) input.ProductUseCase {
    return &productUseCase{
        productRepo: productRepo,
    }
}

func (uc *productUseCase) GetProducts(ctx context.Context, filter *dto.ProductFilter) (*dto.PaginatedOutput[dto.ProductOutput], error) {
    products, err := uc.productRepo.FindAll(ctx, filter)
    if err != nil {
        return nil, err
    }
    
    count, err := uc.productRepo.Count(ctx, filter)
    if err != nil {
        return nil, err
    }
    
    outputs := make([]dto.ProductOutput, len(products))
    for i, p := range products {
        outputs[i] = *uc.toOutput(p)
    }
    
    return dto.NewPaginatedOutput(outputs, filter.Page, filter.Limit, count), nil
}

func (uc *productUseCase) GetProductByID(ctx context.Context, id uint) (*dto.ProductOutput, error) {
    product, err := uc.productRepo.FindByID(ctx, id)
    if err != nil {
        return nil, errors.ErrProductNotFound
    }
    return uc.toOutput(product), nil
}

func (uc *productUseCase) CreateProduct(ctx context.Context, input *dto.ProductInput) (*dto.ProductOutput, error) {
    // Validar entrada
    if err := input.Validate(); err != nil {
        return nil, err
    }
    
    // Criar entidade
    product := entity.NewProduct(
        getValue(input.Name, ""),
        getValue(input.Description, ""),
        getValue(input.Price, 0),
        getValue(input.Stock, 0),
    )
    
    // Validar entidade
    if err := product.Validate(); err != nil {
        return nil, err
    }
    
    // Salvar
    if err := uc.productRepo.Create(ctx, product); err != nil {
        return nil, err
    }
    
    return uc.toOutput(product), nil
}

func (uc *productUseCase) UpdateProduct(ctx context.Context, id uint, input *dto.ProductInput) (*dto.ProductOutput, error) {
    product, err := uc.productRepo.FindByID(ctx, id)
    if err != nil {
        return nil, errors.ErrProductNotFound
    }
    
    // Atualizar campos
    if input.Name != nil {
        product.Name = *input.Name
    }
    if input.Description != nil {
        product.Description = *input.Description
    }
    if input.Price != nil {
        product.Price = *input.Price
    }
    if input.Stock != nil {
        product.Stock = *input.Stock
    }
    if input.Active != nil {
        product.Active = *input.Active
    }
    
    if err := product.Validate(); err != nil {
        return nil, err
    }
    
    if err := uc.productRepo.Update(ctx, product); err != nil {
        return nil, err
    }
    
    return uc.toOutput(product), nil
}

func (uc *productUseCase) DeleteProducts(ctx context.Context, ids []uint) error {
    return uc.productRepo.Delete(ctx, ids)
}

func (uc *productUseCase) toOutput(p *entity.Product) *dto.ProductOutput {
    return &dto.ProductOutput{
        ID:          &p.ID,
        Name:        &p.Name,
        Description: &p.Description,
        Price:       &p.Price,
        Stock:       &p.Stock,
        Active:      &p.Active,
    }
}

func getValue[T any](ptr *T, defaultVal T) T {
    if ptr != nil {
        return *ptr
    }
    return defaultVal
}
```

---

## 🔗 Passo 8: Registrar no DI Container

**Arquivo:** `internal/di/container.go`

```go
// Adicionar no struct Repositories em internal/app/application.go:
type Repositories struct {
    User    output.UserRepository
    Profile output.ProfileRepository
    Product output.ProductRepository  // ← NOVO
}

// Adicionar no struct useCases em internal/di/container.go:
type useCases struct {
    Auth    input.AuthUseCase
    Profile input.ProfileUseCase
    User    input.UserUseCase
    Product input.ProductUseCase  // ← NOVO
}

// Adicionar em initRepositories():
func (c *Container) initRepositories() {
    c.repositories = &app.Repositories{
        User:    repository.NewUserRepository(c.DB),
        Profile: repository.NewProfileRepository(c.DB),
        Product: repository.NewProductRepository(c.DB),  // ← NOVO
    }
}

// Adicionar em initUseCases():
func (c *Container) initUseCases() {
    c.useCases = &useCases{
        // ... existentes ...
        Product: product.NewProductUseCase(c.repositories.Product),  // ← NOVO
    }
}

// Adicionar getter:
func (c *Container) ProductUseCase() input.ProductUseCase {
    return c.useCases.Product
}
```

---

## 🎯 Passo 9: Adicionar à Application

**Arquivo:** `internal/app/application.go`

```go
type Application struct {
    Config  *config.Config
    Log     *logger.Logger
    
    Auth    input.AuthUseCase
    Profile input.ProfileUseCase
    User    input.UserUseCase
    Product input.ProductUseCase  // ← NOVO
    
    Repositories *Repositories
}

// Atualizar o construtor New():
func New(
    cfg *config.Config,
    log *logger.Logger,
    authUC input.AuthUseCase,
    profileUC input.ProfileUseCase,
    userUC input.UserUseCase,
    productUC input.ProductUseCase,  // ← NOVO
    repos *Repositories,
    opts ...Option,
) *Application {
    // ...
}
```

---

## 🌐 Passo 10: Criar o Handler REST

**Arquivo:** `internal/adapter/driver/rest/handler/product_handler.go`

```go
package handler

import (
    "github.com/gofiber/contrib/fiberi18n/v2"
    "github.com/gofiber/fiber/v2"
    
    "github.com/raulaguila/go-api/internal/adapter/driver/rest/middleware"
    "github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
    "github.com/raulaguila/go-api/internal/core/dto"
    "github.com/raulaguila/go-api/internal/core/port/input"
)

type ProductHandler struct {
    useCase     input.ProductUseCase
    handleError func(*fiber.Ctx, error) error
}

func NewProductHandler(router fiber.Router, useCase input.ProductUseCase, accessAuth fiber.Handler) {
    handler := &ProductHandler{
        useCase:     useCase,
        handleError: middleware.DefaultErrorHandler(),
    }
    
    // Middlewares para parser DTOs
    productFilterDTO := middleware.ParseDTO(middleware.DTOConfig{
        ContextKey: localFilter,
        OnLookup:   middleware.Query,
        Model:      &dto.ProductFilter{},
    })
    
    productInputDTO := middleware.ParseDTO(middleware.DTOConfig{
        ContextKey: localDTO,
        OnLookup:   middleware.Body,
        Model:      &dto.ProductInput{},
    })
    
    idParamDTO := middleware.ParseDTO(middleware.DTOConfig{
        ContextKey: localID,
        OnLookup:   middleware.Params,
        Model:      &struct{ ID uint `params:"id"` }{},
    })
    
    idsBodyDTO := middleware.ParseDTO(middleware.DTOConfig{
        ContextKey: localID,
        OnLookup:   middleware.Body,
        Model:      &dto.IDsInput{},
    })
    
    // Rotas protegidas
    router.Use(accessAuth)
    router.Get("", productFilterDTO, handler.getProducts)
    router.Get("/:id", idParamDTO, handler.getProduct)
    router.Post("", productInputDTO, handler.createProduct)
    router.Put("/:id", idParamDTO, productInputDTO, handler.updateProduct)
    router.Delete("", idsBodyDTO, handler.deleteProducts)
}

func (h *ProductHandler) getProducts(c *fiber.Ctx) error {
    filter := c.Locals(localFilter).(*dto.ProductFilter)
    response, err := h.useCase.GetProducts(c.Context(), filter)
    if err != nil {
        return h.handleError(c, err)
    }
    return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ProductHandler) getProduct(c *fiber.Ctx) error {
    id := c.Locals(localID).(*struct{ ID uint `params:"id"` })
    response, err := h.useCase.GetProductByID(c.Context(), id.ID)
    if err != nil {
        return h.handleError(c, err)
    }
    return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ProductHandler) createProduct(c *fiber.Ctx) error {
    input := c.Locals(localDTO).(*dto.ProductInput)
    product, err := h.useCase.CreateProduct(c.Context(), input)
    if err != nil {
        return h.handleError(c, err)
    }
    return presenter.Created(c, fiberi18n.MustLocalize(c, "productCreated"), product)
}

func (h *ProductHandler) updateProduct(c *fiber.Ctx) error {
    id := c.Locals(localID).(*struct{ ID uint `params:"id"` })
    input := c.Locals(localDTO).(*dto.ProductInput)
    product, err := h.useCase.UpdateProduct(c.Context(), id.ID, input)
    if err != nil {
        return h.handleError(c, err)
    }
    return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "productUpdated"), product)
}

func (h *ProductHandler) deleteProducts(c *fiber.Ctx) error {
    ids := c.Locals(localID).(*dto.IDsInput)
    if err := h.useCase.DeleteProducts(c.Context(), ids.IDs); err != nil {
        return h.handleError(c, err)
    }
    return presenter.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "productDeleted"), nil)
}
```

---

## 🛤️ Passo 11: Registrar as Rotas

**Arquivo:** `internal/adapter/driver/rest/server.go`

```go
func (s *Server) setupRoutes() {
    // ... código existente ...
    
    // Registrar handler de produtos
    handler.NewProductHandler(
        s.app.Group("/product"),
        s.appCtx.Product,  // ← Usa o use case da Application
        accessAuth,
    )
    
    // ... resto do código ...
}
```

---

## 🗃️ Passo 12: Criar Migração do Banco (se necessário)

Adicione a tabela no banco de dados:

```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    stock INTEGER DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Ou adicione o modelo à auto-migração no main.go/container.go.

---

## ✅ Checklist Final

```
[ ] Entidade criada (entity/product.go)
[ ] Erros de domínio criados (errors/product_errors.go)
[ ] DTOs criados (dto/product.go)
[ ] Port de entrada criado (port/input/product_usecase.go)
[ ] Port de saída criado (port/output/product_repository.go)
[ ] Modelo GORM criado (model/product_model.go)
[ ] Mapper criado/atualizado (mapper/mapper.go)
[ ] Repository implementado (repository/product_repository.go)
[ ] Use Case implementado (usecase/product/product_usecase.go)
[ ] Container atualizado (di/container.go)
[ ] Application atualizada (app/application.go)
[ ] Handler criado (handler/product_handler.go)
[ ] Rotas registradas (server.go)
[ ] Migração do banco executada
[ ] Build passa (go build ./...)
[ ] Swagger atualizado (swag init)
```

---

## 🧪 Testando

```bash
# Build
go build -o /dev/null ./cmd/backend/main.go

# Atualizar Swagger
swag init -g cmd/backend/main.go -o docs

# Rodar a aplicação
make run

# Testar endpoints
curl http://localhost:9000/product
curl -X POST http://localhost:9000/product -d '{"name":"Produto 1","price":99.90}'
```

---

## 📊 Resumo Visual

```
┌──────────────────────────────────────────────────────────────────┐
│                     NOVO MÓDULO: PRODUCT                         │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  📁 entity/product.go          ← O que é um Produto              │
│  📁 errors/product_errors.go   ← O que pode dar errado           │
│  📁 dto/product.go             ← Dados de entrada/saída          │
│                                                                  │
│  📁 port/input/product_usecase.go    ← Interface do Use Case     │
│  📁 port/output/product_repository.go ← Interface do Repository  │
│                                                                  │
│  📁 usecase/product/product_usecase.go ← Lógica de negócio       │
│                                                                  │
│  📁 model/product_model.go     ← Modelo GORM                     │
│  📁 mapper/mapper.go           ← Conversores                     │
│  📁 repository/product_repository.go ← Acesso ao banco           │
│                                                                  │
│  📁 handler/product_handler.go ← Endpoints REST                  │
│  📁 server.go                  ← Registro das rotas              │
│                                                                  │
│  📁 di/container.go            ← Conecta tudo                    │
│  📁 app/application.go         ← Expõe para interfaces           │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

*Seguindo estes passos, você mantém a arquitetura hexagonal consistente e prepara o código para fácil manutenção e testes.*
