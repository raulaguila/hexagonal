package entity_test

import (
	"testing"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	auth, err := entity.NewAuth(1, true)
	assert.NoError(t, err)

	type args struct {
		name     string
		username string
		email    string
		auth     *entity.Auth
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid User",
			args: args{
				name:     "John Doe",
				username: "johndoe",
				email:    "john@example.com",
				auth:     auth,
			},
			wantErr: false,
		},
		{
			name: "Invalid Name",
			args: args{
				name:     "Jo",
				username: "johndoe",
				email:    "john@example.com",
				auth:     auth,
			},
			wantErr: true,
		},
		{
			name: "Invalid Email",
			args: args{
				name:     "John Doe",
				username: "johndoe",
				email:    "invalid-email",
				auth:     auth,
			},
			wantErr: true,
		},
		{
			name: "Nil Auth",
			args: args{
				name:     "John Doe",
				username: "johndoe",
				email:    "john@example.com",
				auth:     nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entity.NewUser(tt.args.name, tt.args.username, tt.args.email, tt.args.auth)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.args.name, got.Name)
			}
		})
	}
}
