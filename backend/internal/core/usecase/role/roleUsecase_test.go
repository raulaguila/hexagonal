package role_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/usecase/role"
	"github.com/raulaguila/go-api/pkg/apperror"
)

// MockRoleRepo implements output.RoleRepository for testing
type MockRoleRepo struct {
	mock.Mock
}

func (m *MockRoleRepo) Count(ctx context.Context, filter *dto.RoleFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRoleRepo) FindAll(ctx context.Context, filter *dto.RoleFilter) ([]*entity.Role, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entity.Role), args.Error(1)
}

func (m *MockRoleRepo) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepo) Create(ctx context.Context, r *entity.Role) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRoleRepo) Update(ctx context.Context, r *entity.Role) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRoleRepo) Delete(ctx context.Context, ids []uuid.UUID) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func TestGetRoles_Success(t *testing.T) {
	mockRepo := new(MockRoleRepo)
	uc := role.NewRoleUseCase(mockRepo)
	ctx := context.Background()

	roles := []*entity.Role{
		entity.NewRole("Admin", []string{"all"}),
		entity.NewRole("User", []string{"read"}),
	}
	count := int64(2)
	filter := &dto.RoleFilter{
		Filter: dto.Filter{
			Page:  1,
			Limit: 10,
		},
	}

	mockRepo.On("FindAll", ctx, filter).Return(roles, nil)
	mockRepo.On("Count", ctx, filter).Return(count, nil)

	output, err := uc.GetRoles(ctx, filter)

	require.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, uint(count), output.Pagination.TotalItems)
	assert.Len(t, output.Items, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetRoleByID_Success(t *testing.T) {
	mockRepo := new(MockRoleRepo)
	uc := role.NewRoleUseCase(mockRepo)
	ctx := context.Background()

	id := uuid.New()
	existingRole := entity.NewRole("Admin", []string{"all"})
	existingRole.ID = id

	mockRepo.On("FindByID", ctx, id).Return(existingRole, nil)

	output, err := uc.GetRoleByID(ctx, id.String())

	require.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "Admin", *output.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetRoleByID_NotFound(t *testing.T) {
	mockRepo := new(MockRoleRepo)
	uc := role.NewRoleUseCase(mockRepo)
	ctx := context.Background()

	id := uuid.New()

	mockRepo.On("FindByID", ctx, id).Return(nil, apperror.RoleNotFound())

	output, err := uc.GetRoleByID(ctx, id.String())

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.IsType(t, &apperror.Error{}, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateRole_Success(t *testing.T) {
	mockRepo := new(MockRoleRepo)
	uc := role.NewRoleUseCase(mockRepo)
	ctx := context.Background()

	name := "New Role"
	perms := []string{"read", "write"}
	input := &dto.RoleInput{
		Name:        &name,
		Permissions: &perms,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Role")).Return(nil)

	output, err := uc.CreateRole(ctx, input)

	require.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, name, *output.Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRole_Success(t *testing.T) {
	mockRepo := new(MockRoleRepo)
	uc := role.NewRoleUseCase(mockRepo)
	ctx := context.Background()

	id := uuid.New()
	existingRole := entity.NewRole("Old Name", []string{"read"})
	existingRole.ID = id

	newName := "Updated Name"
	input := &dto.RoleInput{
		Name: &newName,
	}

	mockRepo.On("FindByID", ctx, id).Return(existingRole, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(r *entity.Role) bool {
		return r.Name == newName
	})).Return(nil)

	output, err := uc.UpdateRole(ctx, id.String(), input)

	require.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, newName, *output.Name)
	mockRepo.AssertExpectations(t)
}

func TestDeleteRoles_Success(t *testing.T) {
	mockRepo := new(MockRoleRepo)
	uc := role.NewRoleUseCase(mockRepo)
	ctx := context.Background()

	id1 := uuid.New()
	id2 := uuid.New()
	ids := []string{id1.String(), id2.String()}

	mockRepo.On("Delete", ctx, []uuid.UUID{id1, id2}).Return(nil)

	err := uc.DeleteRoles(ctx, ids)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
