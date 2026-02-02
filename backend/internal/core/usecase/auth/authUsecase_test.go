package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/usecase/auth"
	"github.com/raulaguila/go-api/pkg/apperror"
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

func (m *MockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
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

func (m *MockUserRepo) Delete(ctx context.Context, ids []uuid.UUID) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

// MockTokenRepo implements output.TokenRepository for testing
type MockTokenRepo struct {
	mock.Mock
}

func (m *MockTokenRepo) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	args := m.Called(ctx, token, expiration)
	return args.Error(0)
}

func (m *MockTokenRepo) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(bool), args.Error(1)
}

func generateTestKeys() *auth.Config {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	return &auth.Config{
		AccessPrivateKey:  priv,
		AccessExpiration:  15 * time.Minute,
		RefreshPrivateKey: priv,
		RefreshExpiration: 1 * time.Hour,
	}
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	config := generateTestKeys()
	uc := auth.NewAuthUseCase(mockRepo, mockTokenRepo, *config)

	ctx := context.Background()
	username := "testuser"
	password := "password123"

	// Create user with hashed password
	authUser, _ := entity.NewAuth(true)
	authUser.SetPassword(password, time.Now())
	user, _ := entity.NewUser("Test User", username, "test@example.com", authUser)

	mockRepo.On("FindByUsername", ctx, username).Return(user, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil) // Update for token save

	input := &dto.LoginInput{
		Login:    username,
		Password: password,
	}

	output, err := uc.Login(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.AccessToken)
	assert.NotEmpty(t, output.RefreshToken)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	config := generateTestKeys()
	uc := auth.NewAuthUseCase(mockRepo, mockTokenRepo, *config)

	ctx := context.Background()
	username := "nonexistent"

	mockRepo.On("FindByUsername", ctx, username).Return(nil, apperror.UserNotFound())

	input := &dto.LoginInput{
		Login:    username,
		Password: "password",
	}

	output, err := uc.Login(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.IsType(t, &apperror.Error{}, err) // Checks if it returns custom AppError
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	config := generateTestKeys()
	uc := auth.NewAuthUseCase(mockRepo, mockTokenRepo, *config)

	ctx := context.Background()
	username := "testuser"
	password := "password123"

	authUser, _ := entity.NewAuth(true)
	authUser.SetPassword(password, time.Now())
	user, _ := entity.NewUser("Test User", username, "test@example.com", authUser)

	mockRepo.On("FindByUsername", ctx, username).Return(user, nil)

	input := &dto.LoginInput{
		Login:    username,
		Password: "wrongpassword",
	}

	output, err := uc.Login(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	mockRepo.AssertExpectations(t)
}

func TestLogin_DisabledUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	config := generateTestKeys()
	uc := auth.NewAuthUseCase(mockRepo, mockTokenRepo, *config)

	ctx := context.Background()
	username := "disableduser"
	password := "password123"

	authUser, _ := entity.NewAuth(false) // Disabled
	authUser.SetPassword(password, time.Now())
	user, _ := entity.NewUser("Disabled User", username, "disabled@example.com", authUser)

	mockRepo.On("FindByUsername", ctx, username).Return(user, nil)

	input := &dto.LoginInput{
		Login:    username,
		Password: password,
	}

	output, err := uc.Login(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	mockRepo.AssertExpectations(t)
}

func TestMe_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	config := generateTestKeys()
	uc := auth.NewAuthUseCase(mockRepo, mockTokenRepo, *config)

	ctx := context.Background()
	id := uuid.New()

	authUser, err := entity.NewAuth(true) // Valid auth
	require.NoError(t, err)
	user, err := entity.NewUser("Test User", "testuser", "test@example.com", authUser)
	require.NoError(t, err)
	user.ID = id

	mockRepo.On("FindByID", ctx, id).Return(user, nil)

	output, err := uc.Me(ctx, id.String())

	require.NoError(t, err)
	require.NotNil(t, output)
	if output != nil {
		assert.Equal(t, id.String(), *output.ID)
	}
	mockRepo.AssertExpectations(t)
}

func TestMe_InvalidUUID(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	config := generateTestKeys()
	uc := auth.NewAuthUseCase(mockRepo, mockTokenRepo, *config)

	ctx := context.Background()

	output, err := uc.Me(ctx, "invalid-uuid")

	assert.Error(t, err)
	assert.Nil(t, output)
}

func TestLogout_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	config := generateTestKeys()
	uc := auth.NewAuthUseCase(mockRepo, mockTokenRepo, *config)

	ctx := context.Background()
	token := "some-token"

	mockTokenRepo.On("BlacklistToken", ctx, token, config.RefreshExpiration).Return(nil)

	err := uc.Logout(ctx, token)

	assert.NoError(t, err)
	mockTokenRepo.AssertExpectations(t)
}
