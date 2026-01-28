package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/usecase/user"
)

// MockUserRepo implements output.UserRepository for testing
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Count(ctx context.Context, filter *dto.UserFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) FindAll(ctx context.Context, filter *dto.UserFilter) ([]*entity.User, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepo) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepo) FindByToken(ctx context.Context, token string) (*entity.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, u *entity.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepo) Update(ctx context.Context, u *entity.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := user.NewUserUseCase(mockRepo)

	ctx := context.Background()
	name := "John Doe"
	username := "johndoe"
	email := "john@example.com"
	profileID := uint(1)
	status := true

	input := &dto.UserInput{
		Name:      &name,
		Username:  &username,
		Email:     &email,
		ProfileID: &profileID,
		Status:    &status,
	}

	// Expect Create to be called
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	// Expect FindByID to be called for reload (after create)
	// We return a user with ID 0 as it's not set by DB in this mock, but that's fine for flow check
	// We return a user with ID 0 as it's not set by DB in this mock, but that's fine for flow check
	auth, _ := entity.NewAuth(profileID, status)
	expectedUser, _ := entity.NewUser(name, username, email, auth)

	mockRepo.On("FindByID", ctx, uint(0)).Return(expectedUser, nil)

	created, err := uc.CreateUser(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, name, *created.Name)
	mockRepo.AssertExpectations(t)
}
