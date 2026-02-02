package role

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/pkg/apperror"
	"github.com/raulaguila/go-api/pkg/utils"
)

// roleUseCase implements the RoleUseCase interface
type roleUseCase struct {
	roleRepo output.RoleRepository
}

// NewRoleUseCase creates a new RoleUseCase instance
func NewRoleUseCase(roleRepo output.RoleRepository) input.RoleUseCase {
	return &roleUseCase{
		roleRepo: roleRepo,
	}
}

// GetRoles returns a paginated list of roles
func (uc *roleUseCase) GetRoles(ctx context.Context, filter *dto.RoleFilter) (*dto.PaginatedOutput[dto.RoleOutput], error) {
	roles, err := uc.roleRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, apperror.Internal("failed to find roles", err)
	}

	count, err := uc.roleRepo.Count(ctx, filter)
	if err != nil {
		return nil, apperror.Internal("failed to count roles", err)
	}

	return dto.NewPaginatedOutput(
		dto.EntitiesToRoleOutputs(roles, true),
		filter.Page,
		filter.Limit,
		count), nil
}

// ListRoles returns a simple list of roles (id + name)
func (uc *roleUseCase) ListRoles(ctx context.Context, filter *dto.RoleFilter) ([]dto.ItemOutput, error) {
	// Disable pagination for list
	filter.Page = 0
	filter.Limit = 0

	roles, err := uc.roleRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, apperror.Internal("failed to list roles", err)
	}

	outputs := make([]dto.ItemOutput, len(roles))
	for i, role := range roles {
		idStr := role.ID.String()
		outputs[i] = dto.ItemOutput{
			ID:   &idStr,
			Name: &role.Name,
		}
	}

	return outputs, nil
}

// GetRoleByID returns a role by its ID
func (uc *roleUseCase) GetRoleByID(ctx context.Context, id string) (*dto.RoleOutput, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.InvalidInput("id", "invalid uuid format")
	}

	role, err := uc.roleRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, apperror.RoleNotFound()
	}

	return dto.EntityToRoleOutput(role, true), nil
}

// CreateRole creates a new role
func (uc *roleUseCase) CreateRole(ctx context.Context, input *dto.RoleInput) (*dto.RoleOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	role := entity.NewRole(
		utils.Deref(input.Name, ""),
		utils.Deref(input.Permissions, []string{}),
	)
	// No error return from NewRole as seen in previous view_file (it returns *Role)

	if err := role.Validate(); err != nil {
		return nil, err
	}

	if err := uc.roleRepo.Create(ctx, role); err != nil {
		return nil, apperror.Internal("failed to create role", err)
	}

	return dto.EntityToRoleOutput(role, true), nil
}

// UpdateRole updates an existing role
func (uc *roleUseCase) UpdateRole(ctx context.Context, id string, input *dto.RoleInput) (*dto.RoleOutput, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.InvalidInput("id", "invalid uuid format")
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

	role, err := uc.roleRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, apperror.RoleNotFound()
	}

	if input.Name != nil {
		role.UpdateName(*input.Name, time.Now())
	}
	if input.Permissions != nil {
		role.UpdatePermissions(*input.Permissions, time.Now())
	}

	if err := role.Validate(); err != nil {
		return nil, err
	}

	if err := uc.roleRepo.Update(ctx, role); err != nil {
		return nil, apperror.Internal("failed to update role", err)
	}

	return dto.EntityToRoleOutput(role, true), nil
}

// DeleteRoles deletes roles by their IDs
func (uc *roleUseCase) DeleteRoles(ctx context.Context, ids []string) error {
	var uids []uuid.UUID
	for _, id := range ids {
		uid, err := uuid.Parse(id)
		if err != nil {
			return apperror.InvalidInput("ids", "invalid uuid format")
		}
		uids = append(uids, uid)
	}

	if len(uids) == 0 {
		return nil
	}
	return uc.roleRepo.Delete(ctx, uids)
}
